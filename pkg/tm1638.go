// This package provides access to the LEDS and buttons on an LED8KEY board.
// The LED8KEY board uses the TM1638 chip.
//
// See the datasheet for the chip here:
// https://futuranet.it/futurashop/image/catalog/data/Download/TM1638_V1.3_EN.pdf
package pkg

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

/*

  First byte command

  --------------------------------
  01: Set data command

  01_00_I_tMM

  I = 0 for auto incrementing address
      1 for fixed address

  t = 0 for normal mode
      1 for test mode (do not use)

  MM = 00 Write data to the display register
       01 Not allowed
       10 Read key scanning data
	   11 Not allowed

  ---------------------------------
  10: Set display control command

  10_00_D_PPP

  D = 0 for display off
      1 for display on

  PPP  = 000 pulse width 1/16 (brightness)
         001 2/16
		 010 4/16
		 011 10/16
		 100 11/16
		 101 12/16
		 110 13/16
		 111 14/16

  --------------------------------
  11: Set address command

  11_00_AAAA

  AAAA = address 0x00 to 0x0F

  --------------------------------
  Reading data:
    1. strobe goes low
	2. write 1st command byte (read command 0x42:  01_00_0_010)
	3. read byte
	4. read byte
	5. read byte
	6. read byte
	7. stobe goes high

  TODO document the writes

*/

type TM1638 struct {
	STROBE rpio.Pin
	CLK    rpio.Pin
	DIO    rpio.Pin
}

func (x *TM1638) Initialize(pinSTROBE int, pinCLK int, pinDIO int) {
	// Three pins to talk to the board. The DIO pin is for input/output.
	// We'll leave the pin in input mode until we need it (avoids a short
	// in case the board makes drives it unexpectedly).

	x.STROBE = rpio.Pin(pinSTROBE)
	x.CLK = rpio.Pin(pinCLK)
	x.DIO = rpio.Pin(pinDIO)

	x.STROBE.Write(1) // Active low -- start it high
	x.CLK.Write(1)    // Active low -- start it high
	x.DIO.Write(0)    // We'll simulate open-drain

	x.STROBE.Output() // Driven
	x.CLK.Output()    // Driven
	x.DIO.Input()     // We'll simulate open-drain
}

// Create a new LED8Key driver with the given pin numbers. These numbers
// are the RPi's BCM pin numbers -- not the board pin numbers on the IO header.
// See the "What do these numbers mean?" section here: https://pinout.xyz/
func NewTM1638(pinSTROBE int, pinCLK int, pinDIO int) *TM1638 {
	ret := &TM1638{}
	ret.Initialize(pinSTROBE, pinCLK, pinDIO)
	return ret
}

// Twiddle the CLK and DIO lines to send one byte of data.
// Data is sent low-bit first. The chip latches data on the
// falling edge of the clock.
func (x *TM1638) sendByte(value byte) {
	// We use the data-direction here to send the data. We are simulating
	// an open-drain output. The board must have a pullup resistor on DIO.
	// The open-drain prevents both boards from driving the line with
	// opposite values.
	for i := 0; i < 8; i++ {
		if (value & 1) == 1 {
			x.DIO.Input() // Release the line, which is pulled up to "1"
		} else {
			x.DIO.Output() // Drive the line to "0"
		}
		x.CLK.Write(0) // Data is read on high to low transition
		time.Sleep(time.Microsecond)
		value = value >> 1 // Next bit
		time.Sleep(time.Microsecond)
		x.CLK.Write(1) // Get ready for next cycle
		time.Sleep(time.Microsecond)
	}
	x.DIO.Input()
}

// Twiddle the CLK and DIO lines to read one byte of data.
// Data is sent low-bit first. Take the clock low to extract the bit.
// Read the bit just before taking the clock high again.
func (x *TM1638) readByte() byte {
	var ret byte = 0 // Accumulate value here
	x.DIO.Input()    // We are reading
	for i := 0; i < 8; i++ {
		x.CLK.Write(0) // Tell the chip to write its data
		ret = ret << 1 // Move our bits over for the new one
		time.Sleep(time.Millisecond)
		if x.DIO.Read() == rpio.High {
			ret |= 1 // Add in a 1 if the data is 1
		}
		x.CLK.Write(1) // Ready for next cycle
		time.Sleep(time.Microsecond)
	}
	return ret
}

// Configure the brightness of all outputs.
//
// enabled = false to turn the display completely off
// pulseWidth:
//   - 0 =  1/16 (dim)
//   - 1 =  2/16
//   - 2 =  4/16
//   - 3 = 10/16
//   - 4 = 11/16
//   - 5 = 12/16
//   - 6 = 13/16
//   - 7 = 14/16 (bright)
func (x *TM1638) ConfigureDisplay(enabled bool, pulseWidth uint) error {
	// 1. Active strobe
	// 2. Send command
	// 3. Release strobe

	if pulseWidth > 7 {
		return fmt.Errorf("Invalid pulseWidth value: %d", pulseWidth)
	}

	//                     E_ppp
	var cmd byte = 0b10_00_0_000
	pulseWidth = pulseWidth & 7
	if enabled {
		cmd |= 0b00_00_1_000
	}
	cmd |= uint8(pulseWidth)

	x.STROBE.Write(0)
	time.Sleep(time.Microsecond)
	x.sendByte(cmd)
	x.STROBE.Write(1)
	time.Sleep(time.Microsecond)

	return nil
}

// Read up to four bytes of key scanning data.
// Four is all there are.
func (x *TM1638) ReadScanningData(data []byte) error {
	// 1. Active strobe
	// 2. Send read command
	// 3. Read data bytes
	// 4. Release strobe
	if (len(data) < 1) || (len(data) > 4) {
		return fmt.Errorf("Can only read 1 to 4 bytes")
	}
	x.STROBE.Write(0)
	time.Sleep(time.Microsecond)
	x.sendByte(0b01_00_0_010) // Read command
	time.Sleep(time.Microsecond)
	for i := int(0); i < len(data); i++ {
		v := x.readByte()
		data[i] = v
		time.Sleep(time.Microsecond)
	}
	x.STROBE.Write(1)
	time.Sleep(time.Microsecond)
	return nil
}

// Prepare the chip to take data.
//   - autoIncrement = true to bump the address automatically after every write
func (x *TM1638) InitWriteData(autoIncrement bool) error {
	// 1. Active strobe
	// 2. Send the command
	// 3. Release the strobe

	x.STROBE.Write(0)
	time.Sleep(time.Microsecond)
	//                     I
	var cmd byte = 0b01_00_0_000
	if !autoIncrement {
		cmd |= 0b1_000
	}
	x.sendByte(cmd)
	time.Sleep(time.Microsecond)
	x.STROBE.Write(1)
	time.Sleep(time.Microsecond)

	return nil
}

// Send an address followed by stream of bytes.
//   - address = the starting address (0x00 to 0x0F)
//   - data = slice of bytes
func (x *TM1638) WriteData(address uint, data []byte) error {
	// 1. Active strobe
	// 2. Send Address
	// 3. Send each byte of data
	// 4. Release strobe

	if address > 0x0F {
		return fmt.Errorf("Invalid address %d. Must be 0 to 15.", address)
	}
	if len(data) < 1 || len(data) > 16 {
		return fmt.Errorf("Data must be 1 to 16 bytes")
	}
	x.STROBE.Write(0)
	time.Sleep(time.Microsecond)
	address |= uint(0b11_00_0000)
	x.sendByte(byte(address))
	time.Sleep(time.Microsecond)
	for _, v := range data {
		x.sendByte(v)
	}
	x.STROBE.Write(1)
	time.Sleep(time.Microsecond)
	return nil
}

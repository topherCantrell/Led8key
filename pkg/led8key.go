// This package provides access to the LEDS and buttons on an LED8KEY board.
// The LED8KEY board uses the TM1638 chip.
//
// See the datasheet for the chip here:
// https://futuranet.it/futurashop/image/catalog/data/Download/TM1638_V1.3_EN.pdf
package pkg

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
	4. pause
	5. read byte
	6. pause
	7. read byte
	8. pause
	9. read byte
	10. pause
	11. read byte
	12. stobe goes high

  --------------------------------
  Write data (fixed address):
    1. strobe goes low
	2. write 1st command byte (write command 0x48: 01_00_1_000)
	3. strobe goes high
	4. pause
	5. strobe goes low
	6. write address (11_00_AAAA)
	7. pause
	8. write data
	9. pause
	10. strobe goes high
	(can repeat 5-10 address/data for up to 16 total bytes)

  --------------------------------
  Write multi data up to 16 bytes (auto increment address):
    1. strobe goes low
	2. write 1st command byte (write command 0x40: 01_00_0_000)
	3. strobe goes high
	4. pause
	5. strobe goes low
	6. write address (11_00_AAAA)
	7. pause
	8. write data
	9. ... repeat pause+data for up to 16 total bytes
	10. strobe goes high
	(can repeat 6-10)

*/

import (
	// "fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type LED8Key struct {
	STROBE rpio.Pin
	CLK    rpio.Pin
	DIO    rpio.Pin
}

// Create a new LED8Key driver with the given pin numbers. These numbers
// are the RPi's BCM pin numbers -- not the board pin numbers on the IO header.
// See the "What do these numbers mean?" section here: https://pinout.xyz/
func NewLED8Key(pinSTROBE int, pinCLK int, pinDIO int) *LED8Key {

	// Three pins to talk to the board. The DIO pin is for input/output.
	// We'll leave the pin in input mode until we need it (avoids a short
	// in case the board makes drives it unexpectedly).

	ret := &LED8Key{
		STROBE: rpio.Pin(pinSTROBE),
		CLK:    rpio.Pin(pinCLK),
		DIO:    rpio.Pin(pinDIO),
	}

	ret.STROBE.Output()
	ret.CLK.Output()
	ret.DIO.Input()

	return ret
}

// Twiddle the CLK and DIO lines to send one byte of data.
// Data is sent low-bit first. The chip latches data on the
// falling edge of the clock.
func (x *LED8Key) sendByte(value uint8) {
	// We drive the output like open drain:
	// - set the output value to 0
	// - use the direction register as the output signal
}

// Twiddle the CLK and DIO lines to read one byte of data.
// Data is sent low-bit first. Take the clock low to extract the bit.
// Read the bit just before taking the clock high again.
func (x *LED8Key) readByte() uint8 {
	return 0
}

func (x *LED8Key) SetLEDs(value uint8) {
	for i := 0; i < 20; i++ {
		x.STROBE.Toggle()
		time.Sleep(time.Second)
	}
}

// TODO separate this into TM1638 driver + board specific
// board specific
func (x *LED8Key) ReadButtons() uint8 {
	return 0
}

// List of bytes -- one byte per segment
func (x *LED8Key) WriteDisplays() {

}

func (x *LED8Key) WriteNumber(int32) {
	// even signed and floats
}

package pkg

/*
    This is the memory layout for the buttons and LEDs on the
	LED8Key board.

	7 segment bit layout xgfedcba:
	     a
        ---
      f	| | b
	   g---
	  e	| | c
		---
		 d
		    * x

	All of the 16 bytes in the TM1638 chip map to writeable LEDs.

    Memory bytes on chip:
	0    Digit 1 (left most digit)
	1    LED 1 (left most LED) (xxxxxxxL)
	2    Digit 2
	3    LED 2
	4    Digit 3
	5    LED 3
	6    Digit 4
	7    LED 4
	8    Digit 5
	9    LED 5
	10   Digit 6
	11   LED 6
	12   Digit 7
	13   LED 7
	14   Digit 8 (right most digit)
	15   LED 8 (right most LED)

	All four of the read-scanning-keys bytes are used to capture the 8 buttons
	on the board:

	Buttons from left to right, A to H map to bits in the four return bytes:
	0    A000E000
	1    B000F000
	2    C000G000
	3    D000H000
*/

type LED8KEY struct {
	TM1638
	SevenSegFont
	digitBuffer [8]byte
}

func (x *LED8KEY) Initialize(pinSTROBE int, pinCLK int, pinDIO int) {
	x.TM1638.Initialize(pinSTROBE, pinCLK, pinDIO)
	x.ResetFont()
}

// Create a new LED8Key driver with the given pin numbers. These numbers
// are the RPi's BCM pin numbers -- not the board pin numbers on the IO header.
// See the "What do these numbers mean?" section here: https://pinout.xyz/
func NewLED8KEY(pinSTROBE int, pinCLK int, pinDIO int) *LED8KEY {
	ret := &LED8KEY{}
	ret.Initialize(pinSTROBE, pinCLK, pinDIO)
	return ret
}

// Set the status of the LEDs.
// leds = slice of booleans left to right, true means on
func (x *LED8KEY) SetLEDs(leds [8]bool) error {

	for i := 0; i < 8; i++ {
		data := byte(0)
		if leds[i] {
			data = byte(1)
		}
		x.WriteData(i*2+1, []byte{data})
	}

	return nil
}

// Write 8 display digits.
// digits = array of raw bit patterns for each display
func (x *LED8KEY) WriteDigits(digits [8]byte) error {
	// We should consider keeping a back-buffer of all 16 bytes and always
	// write them together with a "Refresh" method.
	for i := 0; i < len(digits); i++ {
		x.digitBuffer[0] = digits[i]
		// Skipping over the LED bytes
		x.WriteData(i*2, x.digitBuffer[0:1])
	}
	return nil
}

// Print the string to the display using the configured font mapping.
// This writes from left to right and blanks any unused digits to the right.
// chars = the text string.
func (x *LED8KEY) WriteString(chars string) error {
	err := x.BuildDigits(chars, 8, x.digitBuffer[:])
	if err != nil {
		return err
	}

	return x.WriteDigits(x.digitBuffer)
}

// Read the 8 buttons
// Returns an array of booleans from left to right, true means pressed
func (x *LED8KEY) ReadButtons(buttons *[8]bool) error {

	data := []byte{0, 0, 0, 0}
	err := x.ReadScanningData(data)
	if err != nil {
		return err
	}

	buttons[0] = data[0]&0x80 > 0
	buttons[1] = data[1]&0x80 > 0
	buttons[2] = data[2]&0x80 > 0
	buttons[3] = data[3]&0x80 > 0
	buttons[4] = data[0]&0x08 > 0
	buttons[5] = data[1]&0x08 > 0
	buttons[6] = data[2]&0x08 > 0
	buttons[7] = data[3]&0x08 > 0

	return nil
}

package pkg

import (
	"fmt"
)

/*
    This is the memory layout for the buttons and LEDs on the
	Disp16Key board.

		7 segment bit layout:
	     a
        ---
      f	| | b
	   g---
	  e	| | c
		---
		 d
		    * x

	Only 8 bytes in the TM1638 chip map to LED segments.

    Memory bytes on chip:
	0    all "a" segments
	1    unused
	2    all "b" segments
	3    unused
	4    all "c" segments
	5    unused
	6    all "d" segments
	7    unused
	8    all "e" segments
	9    unused
	10   all "f" segments
	11   unused
	12   all "g" segments
	13   unused
	14   all "x" segments
	15   unused

	In each byte:
	bit 0: far right digit
	bit 1:
	bit 2:
	bit 3:
	bit 4:
	bit 5:
	bit 6:
	bit 7: far left digit

	All four of the read-scanning-keys bytes are used to capture the 16 buttons
	on the board:

	// Returns from pressing each button one at a time
	//     Column 1       Column 2       Column 3       Column 4
	//  [20,00,00,00]  [02,00,00,00]  [00,20,00,00]  [00,02,00,00] Row 1
	//  [00,00,20,00]  [00,00,02,00]  [00,00,00,20]  [00,00,00,02] Row 2
	//  [40,00,00,00]  [04,00,00,00]  [00,40,00,00]  [00,04,00,00] Row 3
	//  [00,00,40,00]  [00,00,04,00]  [00,00,00,40]  [00,00,00,04] Row 4

*/

type DISP16KEY struct {
	TM1638
	SevenSegFont
	digitBuffer [8]byte
}

// Create a new DISP16KEY driver with the given pin numbers. These numbers
// are the RPi's BCM pin numbers -- not the board pin numbers on the IO header.
// See the "What do these numbers mean?" section here: https://pinout.xyz/
func NewDISP16KEY(pinSTROBE GPIOPin, pinCLK GPIOPin, pinDIO GPIOPin) *DISP16KEY {
	ret := &DISP16KEY{}
	ret.TM1638 = *NewTM1638(pinSTROBE, pinCLK, pinDIO)
	ret.ResetFont()
	return ret
}

// Write 8 display digits.
// digits = array of raw bit patterns for each display
func (x *DISP16KEY) WriteDigits(digits [8]byte) error {

	// For the 16-key, we have to convert the digits into a different format than used by the 8-key (and thus our other processing methods as well)
	digits = convertEightKeyDigits(digits)
	fmt.Println(">:>", digits)

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
func (x *DISP16KEY) WriteString(chars string) error {
	err := x.BuildDigits(chars, 8, x.digitBuffer[:])
	if err != nil {
		return err
	}

	return x.WriteDigits(x.digitBuffer)
}

// Read the 16 buttons
// Fills in button booleans from left to right and top to bottom
func (x *DISP16KEY) ReadButtons(buttons *[16]bool) error {

	data := []byte{0, 0, 0, 0}
	err := x.ReadScanningData(data)
	if err != nil {
		return err
	}

	buttons[0] = data[0]&0x20 > 0
	buttons[1] = data[0]&0x02 > 0
	buttons[2] = data[1]&0x20 > 0
	buttons[3] = data[1]&0x02 > 0
	buttons[4] = data[2]&0x20 > 0
	buttons[5] = data[2]&0x02 > 0
	buttons[6] = data[3]&0x20 > 0
	buttons[7] = data[3]&0x02 > 0
	buttons[8] = data[0]&0x40 > 0
	buttons[9] = data[0]&0x04 > 0
	buttons[10] = data[1]&0x40 > 0
	buttons[11] = data[1]&0x04 > 0
	buttons[12] = data[2]&0x40 > 0
	buttons[13] = data[2]&0x04 > 0
	buttons[14] = data[3]&0x40 > 0
	buttons[15] = data[3]&0x04 > 0

	return nil
}

// Converts display bytes formatted for the 8-key to the format that is used by the 16-key
// See the top of this file for details on the 16-key data format
func convertEightKeyDigits(digits [8]byte) [8]byte {
	var outputDigits [8]byte

	for i, digit := range digits {

		byteIndex := 0

		// For each bit in 'digit'
		for bitMask := byte(1); bitMask > 0; bitMask = bitMask << 1 {
			// If we have a high bit at the bit mask position
			if bitMask&digit == bitMask {
				outputDigits[byteIndex] |= (128 >> i)
			}

			byteIndex++
		}
	}

	return outputDigits
}

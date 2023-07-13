package pkg

import (
	"fmt"
)

/*
    This is the memory layout for the buttons and LEDs on the
	LED8Key board.

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

	7 segment bit layout xgfedcba:
	     a
        ---
      f	| | b
	   g---
	  e	| | c
		---
		 d
		    * x

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
	font map[int]byte
}

func (x *LED8KEY) Initialize(pinSTROBE int, pinCLK int, pinDIO int) {
	x.TM1638.Initialize(pinSTROBE, pinCLK, pinDIO)

	// Limited font mapping from string characters to bit patterns. The user
	// can extend/change this mapping as needed.

	x.font = map[int]byte{
		' ': 0b0_0000000,
		'0': 0b0_0111111,
		'1': 0b0_0000110,
		'2': 0b0_1011011,
		'3': 0b0_1001111,
		'4': 0b0_1100110,
		'5': 0b0_1101101,
		'6': 0b0_1111101,
		'7': 0b0_0000111,
		'8': 0b0_1111111,
		'9': 0b0_1101111,
		'.': 0b1_0000000, // The PrintString will attempt to combine
		'-': 0b0_1000000, // Minus sign
		// Useful for hex
		'A': 0b0_1110111,
		'B': 0b0_1111100,
		'C': 0b0_0111001,
		'D': 0b0_1011110,
		'E': 0b0_1111001,
		'F': 0b0_1110001,
		// Some random letters for example
		'H': 0b0_1110110,
		'i': 0b0_0000100,
		'L': 0b0_0111000,
		'o': 0b0_1011100,
	}
}

// Create a new LED8Key driver with the given pin numbers. These numbers
// are the RPi's BCM pin numbers -- not the board pin numbers on the IO header.
// See the "What do these numbers mean?" section here: https://pinout.xyz/
func NewLED8KEY(pinSTROBE int, pinCLK int, pinDIO int) *LED8KEY {
	ret := &LED8KEY{}
	ret.Initialize(pinSTROBE, pinCLK, pinDIO)
	return ret
}

// Fill the 16 bytes of memory with a given value.
func (x *LED8KEY) FillDisplay(fillValue byte) error {
	for i := 0; i < 16; i++ {
		x.tempBuffer[i] = fillValue
	}
	err := x.WriteData(0, x.tempBuffer)
	if err != nil {
		return err
	}
	return nil
}

// Set the status of the LEDs.
// leds = slice of booleans left to right, true means on
func (x *LED8KEY) SetLEDs(leds []bool) error {
	if len(leds) != 8 {
		return fmt.Errorf("leds slice must be 8 booleans")
	}

	for i := 0; i < 8; i++ {
		data := byte(0)
		if leds[i] {
			data = byte(1)
		}
		x.WriteData(i*2+1, []byte{data})
	}

	return nil
}

// Write 1 to 16 display digits.
// start = the first digit to begin drawing.
// digits = slice of raw bit patterns for each display.
func (x *LED8KEY) WriteDigits(start int, digits []byte) error {
	if len(digits) > 8 {
		return fmt.Errorf("No more than 8 digits")
	}
	if start < 0 || start > (8-len(digits)) {
		return fmt.Errorf("Invalid start position %d with length of %d", start, len(digits))
	}
	for i := 0; i < len(digits); i++ {
		x.tempBuffer[0] = digits[i]
		x.WriteData(start*2+i*2, x.tempBuffer[0:1])
	}
	return nil
}

// Print the string to the display.
// This method maps characters to segment bit patterns. The user can extend/change
// this font mapping.
// start = the first digit to begin printing.
// chars = the text string.
func (x *LED8KEY) WriteString(start int, chars string) error {

	previous := -1 // No previous-position yet
	pos := 0       // Next digit to fill

	for i := 0; i < len(chars); i++ {
		if chars[i] == '.' {
			// If this is a period, we'll try to merge it with the previous digit
			if previous >= 0 {
				x.tempBuffer[previous] |= 0b1_0000000
				previous = -1
				continue // No new digit ... continue with next character
			}
		}
		// Lookup the segment bit pattern
		value, exist := x.font[int(chars[i])]
		if !exist {
			return fmt.Errorf("No font mapping for '%c' in '%s'.", chars[i], chars)
		}
		if i > 15 {
			return fmt.Errorf("Maximum of 16 digits")
		}
		// Add the value
		x.tempBuffer[pos] = value
		pos++
		previous = i
		if chars[i] == '.' {
			// Decimal points cannot be merged to dots
			previous = -1
		}

	}

	return x.WriteDigits(start, x.tempBuffer[0:pos])
}

// Get the current font mapping for PrintString. Mutate this map as needed.
func (x *LED8KEY) GetMutableFont() map[int]byte {
	return x.font
}

// Read the 8 buttons
// Returns an array of booleans from left to right, true means pressed
func (x *LED8KEY) ReadButtons(buttons []bool) error {

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

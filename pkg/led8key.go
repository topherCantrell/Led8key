package pkg

import (
	"fmt"
)

/*
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
*/

type LED8KEY struct {
	TM1638
	font map[int]byte
}

func (x *LED8KEY) Initialize(pinSTROBE int, pinCLK int, pinDIO int) {
	x.TM1638.Initialize(pinSTROBE, pinCLK, pinDIO)
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
		'-': 0b0_1000000,
		// Useful for hex
		'A': 0b0_1110111,
		'B': 0b0_1111100,
		'C': 0b0_0111001,
		'D': 0b0_1011110,
		'E': 0b0_1111001,
		'F': 0b0_1110001,
		// Some random letters
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

func (x *LED8KEY) FillDisplay(fillValue byte) error {
	// todo there must be better syntax
	data := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < 16; i++ {
		data[i] = fillValue
	}
	err := x.WriteData(0, data)
	if err != nil {
		return err
	}
	return nil
}

func (x *LED8KEY) SetLEDs(leds []bool) error {
	if len(leds) != 8 {
		return fmt.Errorf("leds slice must be 8 booleans")
	}

	for i := 0; i < 8; i++ {
		data := byte(0)
		if leds[i] {
			data = byte(1)
		}
		x.WriteData(uint(i*2+1), []byte{data})
	}

	return nil
}

func (x *LED8KEY) WriteDigits(start int, segments []byte) error {
	if len(segments) > 8 {
		return fmt.Errorf("No more than 8 digits")
	}
	if start < 0 || start > (8-len(segments)) {
		return fmt.Errorf("Invalid start position")
	}
	// TODO there must be a better way than this
	data := []byte{0}
	for i := 0; i < len(segments); i++ {
		data[0] = segments[i]
		// Is this common? all this casting the sign away? Maybe just use an int?
		// But it is nice to have the compiler reject negatives
		x.WriteData(uint(start*2+i*2), data)
	}
	return nil
}

func (x *LED8KEY) WriteString(start int, chars string) error {
	// TODO handle the decimal point by combining it to the left
	// if there is a digit to append. Otherwise use space.
	data := make([]byte, len(chars))
	for i := 0; i < len(chars); i++ {
		value, exist := x.font[int(chars[i])]
		if !exist {
			return fmt.Errorf("No font mapping for '%c' in '%s'.", chars[i], chars)
		}
		data[i] = value
	}
	return x.WriteDigits(start, data)
}

func (x *LED8KEY) GetMutableFont() map[int]byte {
	return x.font
}

func (x *LED8KEY) ReadButtons() ([]bool, error) {
	// Buttons from left to right, A to H map to bits in the return:
	//   0    A000E000
	//   1    B000F000
	//   2    C000G000
	//   3    D000H000
	data := []byte{0, 0, 0, 0}
	err := x.ReadScanningData(data)
	if err != nil {
		return nil, err
	}
	buttons := []bool{false, false, false, false, false, false, false, false}

	if data[0]&0x80 > 0 {
		buttons[0] = true
	}
	if data[1]&0x80 > 0 {
		buttons[1] = true
	}
	if data[2]&0x80 > 0 {
		buttons[2] = true
	}
	if data[3]&0x80 > 0 {
		buttons[3] = true
	}

	if data[0]&0x08 > 0 {
		buttons[4] = true
	}
	if data[1]&0x08 > 0 {
		buttons[5] = true
	}
	if data[2]&0x08 > 0 {
		buttons[6] = true
	}
	if data[3]&0x08 > 0 {
		buttons[7] = true
	}

	return buttons, nil
}

package pkg

import "fmt"

type SevenSegFont struct {
	font map[int]byte
}

func NewSevenSegFont() *SevenSegFont {
	ret := &SevenSegFont{}
	ret.ResetFont()
	return ret
}

// Reset the font back to the default mapping
func (x *SevenSegFont) ResetFont() {
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

// Get the current font mapping for PrintString. Mutate this map as needed.
func (x *SevenSegFont) GetMutableFont() map[int]byte {
	return x.font
}

// Build a slice of digit values for display.
// This method maps characters to segment bit patterns. The user can extend/change
// this font mapping. This method merges decimal points into the digits.
// chars = the text string
// maxDigits = the maximum digits to create
// outDigits = the returned slice of digit data
func (x *SevenSegFont) BuildDigits(chars string, numDigits int, outDigits []byte) error {

	previous := -1 // No previous-position yet
	pos := 0       // Next digit to fill

	for i := 0; i < len(chars); i++ {
		if chars[i] == '.' {
			// If this is a period, we'll try to merge it with the previous digit
			if previous >= 0 {
				outDigits[previous] |= 0b1_0000000
				previous = -1
				continue // No new digit ... continue with next character
			}
		}
		// Lookup the segment bit pattern
		value, exist := x.font[int(chars[i])]
		if !exist {
			return fmt.Errorf("No font mapping for '%c' in '%s'.", chars[i], chars)
		}
		if i > numDigits {
			return fmt.Errorf("Exceeded number of %d digits", numDigits)
		}
		// Add the value
		outDigits[pos] = value
		pos++
		previous = i
		if chars[i] == '.' {
			// Decimal points cannot be merged to dots
			previous = -1
		}

	}

	for i := pos; i < numDigits; i++ {
		// Blank untouched digits
		outDigits[i] = 0
	}

	return nil
}

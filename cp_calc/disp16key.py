from sevensegfont import SevenSegFont
from tm1638 import TM1638

class Disp16Key(TM1638):
    """
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

    """

    def __init__(self, strobe_pin, clock_pin, dio_pin):
        super().__init__(strobe_pin, clock_pin, dio_pin)
        self._font = SevenSegFont()
        self._digit_buffer = [0,0,0,0,0,0,0,0]
        self._scan_buffer = [0,0,0,0]
        
    def write_string(self, chars):
        digits = self._font.build_digits(chars)
        while len(digits)<8:
            digits.append(0)
        digits = digits[0:8]
        self.write_digits(digits) 
    
    def write_digits(self, digits):
        """Pass exactly 8 digits
        
        We would like to have a single byte hold the leds for a single digit. But instead, this
        hardware is mapped with the leds for a single digit spread out over 8 bytes. We do the
        math here to unscramble the hardware.        
        """
        dig2 = self._convert_8_key_digits(digits)
        for i in range(8):
            self.write_data(i*2,dig2[i:i+1])        

    def read_buttons(self, buttons):        
        data = self._scan_buffer     
        self.read_scanning_data(data)           
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

    def _convert_8_key_digits(self, digits):
        data = self._digit_buffer
        i = 0
        while i<8:
            data[i] = 0
            i += 1
            
        i = 0
        while i<8:        
            digit = digits[i]
            byte_index = 0
            bit_mask = 1
            while bit_mask<256:
                if (bit_mask & digit) == bit_mask:                    
                    self._digit_buffer[byte_index] = self._digit_buffer[byte_index] | (128 >> i)
                bit_mask = bit_mask << 1
                byte_index += 1
            i += 1
        return self._digit_buffer
        
    

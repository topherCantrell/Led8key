
class SevenSegFont:

    def __init__(self):
        self.font = {
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
            '.': 0b1_0000000, # The PrintString will attempt to combine
            '-': 0b0_1000000, # Minus sign
            # Useful for hex
            'A': 0b0_1110111,
            'B': 0b0_1111100,
            'C': 0b0_0111001,
            'D': 0b0_1011110,
            'E': 0b0_1111001,
            'F': 0b0_1110001,
            # Some random letters for example
            # h todo
            # I
            # l
            # O
            'H': 0b0_1110110,
            'i': 0b0_0000100,
            'L': 0b0_0111000,
            'o': 0b0_1011100,
        }

    def build_digits(self, chars):
        """ Build a list of digit values for display.
        This method maps characters to segment bit patterns. The user can extend/change
        this font mapping. This method merges decimal points into the digits.        
        chars = the text string
        maxDigits = the maximum digits to create
        outDigits = the returned slice of digit data
        """

        previous = False        
        ret = []

        for i in range(len(chars)):
            if chars[i] == '.':
                # If this is a period, we'll try to merge it with the previous digit
                if previous:
                    ret[-1] |= 0b1_0000000
                    previous = False
                    continue # No new digit ... continue with next character

            # Lookup the segment bit pattern
            value = self.font[chars[i]]
            ret.append(value)
            previous = True
            if chars[i] == '.':
                # Decimal points cannont be merged to decimal points
                previous = False
            
        return ret
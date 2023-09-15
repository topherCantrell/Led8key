import digitalio
import board
import time


class TM1638:

    #   First byte command
    #
    #   --------------------------------
    #   01: Set data command
    #
    #   01_00_I_tMM
    #
    #   I = 0 for auto incrementing address
    #       1 for fixed address
    #
    #   t = 0 for normal mode
    #       1 for test mode (do not use)
    #
    #   MM = 00 Write data to the display register
    #        01 Not allowed
    #        10 Read key scanning data
    # 	   11 Not allowed
    #
    #   ---------------------------------
    #   10: Set display control command
    #
    #   10_00_D_PPP
    #
    #   D = 0 for display off
    #       1 for display on
    #
    #   PPP  = 000 pulse width 1/16 (brightness)
    #          001 2/16
    # 		 010 4/16
    # 		 011 10/16
    # 		 100 11/16
    # 		 101 12/16
    # 		 110 13/16
    # 		 111 14/16
    #
    #   --------------------------------
    #   11: Set address command
    #
    #   11_00_AAAA
    #
    #   AAAA = address 0x00 to 0x0F
    #
    #   --------------------------------
    #   Reading data:
    #     1. strobe goes low
    # 	2. write 1st command byte (read command 0x42:  01_00_0_010)
    # 	3. read byte
    # 	4. read byte
    # 	5. read byte
    # 	6. read byte
    # 	7. stobe goes high
    #
    #   TODO document the writes

    def __init__(self, strobe_pin, clock_pin, dio_pin):
        strobe_pin.direction = digitalio.Direction.OUTPUT
        strobe_pin.value = True
        self._strobe_pin = strobe_pin

        clock_pin.direction = digitalio.Direction.OUTPUT
        clock_pin.value = True
        self._clock_pin = clock_pin

        dio_pin.direction = digitalio.Direction.OUTPUT
        dio_pin.drive_mode = digitalio.DriveMode.OPEN_DRAIN
        dio_pin.value = True
        self._dio_pin = dio_pin

    def send_byte(self, value):
        # Remember, the DIO pin is in OPEN_DRAIN. The library manipulates
        # it correctly with "value". The chip latches data on the
        # falling edge of the clock.
        for i in range(8):
            self._dio_pin.value = (value & 1)
            self._clock_pin.value = False
            time.sleep(0.000_001)
            self._clock_pin.value = True
            time.sleep(0.000_001)
            value = value >> 1                    

    def read_byte(self):
        # Data is sent low-bit first. Take the clock low to extract the bit.
        # Read the bit just before taking the clock high again.
        self._dio_pin.value = True  # Let the TM1638 drive the line
        ret = 0
        for i in range(8):
            self._clock_pin.value = False
            ret = ret << 1
            time.sleep(0.000_001)
            if self._dio_pin.value:
                ret |= 1
            self._clock_pin.value = True
            time.sleep(0.000_001)
        # We leave the clock high (idle)
        return ret
    
    def configure_display(self, enabled, pulse_width):
        """Configure the brightness of all outputs
        
        enabled = false to turn the display completely off
        pulseWidth:
          - 0 =  1/16 (dim)
          - 1 =  2/16
          - 2 =  4/16
          - 3 = 10/16
          - 4 = 11/16
          - 5 = 12/16
          - 6 = 13/16
          - 7 = 14/16 (bright)
        """
        
        cmd = 0b10_00_0_000
        if enabled:
            cmd |= 0b00_00_1_000

        pulse_width = pulse_width & 7
        cmd |= pulse_width

        self._strobe_pin.value = False
        time.sleep(0.000_001)
        self.send_byte(cmd)
        self._strobe_pin.value = True
        time.sleep(0.000_001)

    def read_scanning_data(self, buffer):
        """Read up to 4 bytes of scanning data
        
        This method fills out the list you provide. Key scanning happens a lot --
        no need to thrash memory building a new buffer every time.
        """
        self._strobe_pin.value = False
        time.sleep(0.000_001)
        self.send_byte(0b01_00_0_010)  # Read command
        time.sleep(0.000_001)
        for i in range(len(buffer)):
            buffer[i] = self.read_byte()        
            time.sleep(0.000_001)        
        self._strobe_pin.value = True
        time.sleep(0.000_001)

    def init_write_data(self, auto_increment=True):
        """Init the display for writing
        
        Once it is ready, you can keep doing writes. After a read, though, you
        need to init-write again.
        """
        self._strobe_pin.value = False
        time.sleep(0.000_001)
        cmd = 0b01_00_0_000
        if not auto_increment:
            cmd |= 0b1_000
        self.send_byte(cmd)
        time.sleep(0.000_001)
        self._strobe_pin.value = True
        time.sleep(0.000_001)

    def write_data(self, address, data):        
        """Write the bytes to the display

        You need to init the display once before a series of writes.
        """
        self._strobe_pin.value = False
        time.sleep(0.000_001)
        address &= 0x0F
        address |= 0b11_00_0000
        self.send_byte(address)
        time.sleep(0.000_001)
        for v in data:        
            self.send_byte(v)
        self._strobe_pin.value = True
        time.sleep(0.000_001)

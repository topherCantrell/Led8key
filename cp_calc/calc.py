import board
import digitalio
import time
from disp16key import Disp16Key


class Calculator:

    def __init__(self, disp):
        self._bts = [False]*16
        self._disp = disp
        self._disp.configure_display(enabled=True, pulse_width=7)
        self._disp.init_write_data()

    def get_calc_key(self):
        while True:
            self._disp.read_buttons(self._bts)
            print('>',self._bts)
            time.sleep(1)
            if not True in self._bts:
                break

        while True:
            self._disp.read_buttons(self._bts)
            time.sleep(1)
            if True in self._bts:
                break

        print(self._bts)
        if self._bts[12] and self._bts[14]:
            return 16
                    
        return self._bts.index(True)

    

 

    
strobe_pin = digitalio.DigitalInOut(board.GP28)
dio_pin = digitalio.DigitalInOut(board.GP27)
clock_pin = digitalio.DigitalInOut(board.GP26)

disp = Disp16Key(strobe_pin, clock_pin, dio_pin)
c = Calculator(disp)

c._disp.write_string('8')

time.sleep(10000)

while True:
    k = c.get_calc_key()
    c._disp.init_write_data()
    c._disp.write_string(str(k))
    print('updated')

    print('S L E E P I N G')
    time.sleep(10)
    print('awake')

    

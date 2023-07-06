// This package provides access to the LEDS and buttons on an LED8KEY board.
// The LED8KEY board uses the TM1638 chip.
//
// See the datasheet for the chip here:
// https://futuranet.it/futurashop/image/catalog/data/Download/TM1638_V1.3_EN.pdf
package pkg

// Driver struct with pin defs passed into new
// private bit-bangers
// public access methods
//   - set LEDs
//   - read buttons
//   - set display list of list of booleans -- 8 booleans per digit, 8 digits
//   - draw number 10_000_000

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

func (x *LED8Key) SetLEDs(value uint8) {
	for i := 0; i < 20; i++ {
		x.STROBE.Toggle()
		time.Sleep(time.Second)
	}
}

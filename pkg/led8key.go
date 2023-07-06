package pkg

import (
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

// Driver struct with pin defs passed into new
// private bit-bangers
// public access methods
//   - set LEDs
//   - read buttons
//   - set display list of list of booleans -- 8 booleans per digit, 8 digits
//   - draw number 10_000_000

type LED8Key struct {
	pinSTROBE int
	pinCLK    int
	pinDIO    int
}

func DDinit() {
	// This should happen once when the module is imported. TODO actually
	// it should happen ONCE no matter how many rpio users are imported.
	// Maybe open and close goes with the main function.
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}
}

func NewLED8Key(pinSTROBE int, pinCLK int, pinDIO int) *LED8Key {
	ret := &LED8Key{
		pinSTROBE: pinSTROBE,
		pinCLK:    pinCLK,
		pinDIO:    pinDIO,
	}

	return ret
}

func (x *LED8Key) SayHi() {
	fmt.Println("I AM HERE")
}

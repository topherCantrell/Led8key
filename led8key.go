package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

// Driver struct with pin defs passed into new
// private bit-bangers
// public access methods
//   - set LEDs
//   - read buttons
//   - set display list of list of booleans -- 8 booleans per digit, 8 digits
//   - draw number 10_000_000

func main() {
	fmt.Println("opening gpio")
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	pin := rpio.Pin(18)
	pin.Output()

	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second)
	}
}

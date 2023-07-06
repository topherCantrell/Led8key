package main

import (
	"fmt"

	"github.com/stianeikeland/go-rpio"
	"github.com/topherCantrell/go-led8key/pkg"
)

func main() {

	// Open the rpio once for all using packages (right now just go-led8key)
	fmt.Println("opening gpio")
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	p := pkg.NewLED8Key(1, 2, 3)
	p.SetLEDs(0b10101010)

}

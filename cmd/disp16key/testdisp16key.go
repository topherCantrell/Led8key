package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/topherCantrell/go-led8key/pkg"
)

func main() {

	// Open the rpio once for all using packages (right now just go-led8key)
	// fmt.Println("opening gpio")
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	p := pkg.NewDISP16KEY(17, 27, 22)

	err = p.ConfigureDisplay(true, 7)
	if err != nil {
		fmt.Println("ConfigureDisplay:", err)
	}

	time.Sleep(time.Second)

	p.WriteString("3.14314")

	buttons := [16]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false}

	p.ReadButtons(&buttons)

	fmt.Println(buttons)

}

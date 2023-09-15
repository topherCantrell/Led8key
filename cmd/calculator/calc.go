package main

import (
	"fmt"
	"machine"
	"time"

	"github.com/topherCantrell/go-led8key/pkg"
)

type TinyGoGPIOPin struct {
	TGpin machine.Pin
}

func (x TinyGoGPIOPin) Write(state bool) {
	if state {
		x.TGpin.High()
	} else {
		x.TGpin.Low()
	}
}

func (x TinyGoGPIOPin) Read() bool {
	return x.TGpin.Get()
}

func (x TinyGoGPIOPin) Input() {
	x.TGpin.Configure(machine.PinConfig{Mode: machine.PinInput})
}

func (x TinyGoGPIOPin) Output() {
	x.TGpin.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func main() {

	strobe := TinyGoGPIOPin{machine.GPIO28}
	clk := TinyGoGPIOPin{TGpin: machine.GPIO26}
	dio := TinyGoGPIOPin{TGpin: machine.GPIO27}

	p := pkg.NewDISP16KEY(strobe, clk, dio)

	err := p.ConfigureDisplay(true, 7)
	if err != nil {
		fmt.Println("ConfigureDisplay:", err)
	}

	time.Sleep(time.Second)

	p.WriteString("9.87654")

	// strobe.Output()
	// for {
	// 	strobe.Write(true)
	// 	time.Sleep(time.Second * 5)
	// 	strobe.Write(false)
	// 	time.Sleep(time.Second * 5)
	// }

	// buttons := [16]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false}

	// p.ReadButtons(&buttons)

	// fmt.Println(buttons)

}

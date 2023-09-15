package pkg

import "github.com/stianeikeland/go-rpio"

type RPiGPIOPin struct {
	RpiPin rpio.Pin
}

func (x RPiGPIOPin) Write(state bool) {
	if state {
		x.RpiPin.Write(1)
	} else {
		x.RpiPin.Write(0)
	}
}

func (x RPiGPIOPin) Read() bool {
	return x.RpiPin.Read() == rpio.High
}

func (x RPiGPIOPin) Input() {
	x.RpiPin.Input()
}

func (x RPiGPIOPin) Output() {
	x.RpiPin.Output()
}

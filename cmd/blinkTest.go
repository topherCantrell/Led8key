package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/topherCantrell/go-led8key/pkg"
)

func testFlash(p *pkg.LED8KEY) {

	for i := 0; i < 10; i++ {
		err := p.ConfigureDisplay(true, 7)
		if err != nil {
			fmt.Println("ConfigureDisplay:", err)
		}
		time.Sleep(time.Second)
		err = p.ConfigureDisplay(true, 1)
		if err != nil {
			fmt.Println("ConfigureDisplay:", err)
		}
		time.Sleep(time.Second)
	}

}

func testButtons(p *pkg.LED8KEY) {

	err := p.InitWriteData(true)
	if err != nil {
		fmt.Println("InitWriteData:", err)
	}
	time.Sleep(time.Second)

	buttons := make([]bool, 8)

	for {
		err := p.ReadButtons(buttons)
		if err != nil {
			fmt.Println("ReadButtons:", err)
		}
		fmt.Println(">>>", buttons)
		time.Sleep(time.Second)
	}

}

func testFill(p *pkg.LED8KEY) {

	err := p.InitWriteData(true)
	if err != nil {
		fmt.Println("InitWriteData:", err)
	}
	time.Sleep(time.Second)

	for {
		err := p.FillDisplay(0x00)
		if err != nil {
			fmt.Println("ConfigureDisplay:", err)
		}

		time.Sleep(time.Second)

		err = p.FillDisplay(0xFF)
		if err != nil {
			fmt.Println("ConfigureDisplay:", err)
		}

		time.Sleep(time.Second)
	}
}

func testHello(p *pkg.LED8KEY) {

	bufferA := []byte{
		0b01110110, // Digit 1 (left most digit)
		0x00,       // LED 1 (left most LED) (xxxxxxxL)
		0b01111001, // Digit 2
		0x00,       // LED 2
		0b00111000, // Digit 3
		0x00,       // LED 3
		0b00111000, // Digit 4
		0x00,       // LED 4
		0b00111111, // Digit 5
		0x00,       // LED 5
		0b10000010, // Digit 6
		0x00,       // LED 6
		0b10000010, // Digit 7
		0x00,       // LED 7
		0b10000010, // Digit 8 (right most digit)
		0x00,       // LED 8 (right most LED)
	}

	err := p.InitWriteData(true)
	if err != nil {
		fmt.Println("InitWriteData:", err)
	}
	time.Sleep(time.Second)

	for {
		err = p.WriteData(0, bufferA)
		if err != nil {
			fmt.Println("WriteData:", err)
		}
		time.Sleep(time.Second)

		bufferB := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
		err = p.WriteData(0, bufferB)
		if err != nil {
			fmt.Println("WriteData:", err)
		}
		time.Sleep(time.Second)

		bufferC := []byte{0xAA, 1, 0x55, 0, 0xAA, 1, 0x55, 0, 0xAA, 1, 0x55, 0, 0xAA, 1, 0x55, 0}
		err = p.WriteData(0, bufferC)
		if err != nil {
			fmt.Println("WriteData:", err)
		}
		time.Sleep(time.Second)
	}

}

func testLEDs(p *pkg.LED8KEY) {
	err := p.InitWriteData(true)
	if err != nil {
		fmt.Println("InitWriteData:", err)
	}
	time.Sleep(time.Second)

	for {
		p.SetLEDs([]bool{true, false, true, false, true, false, true, false})
		time.Sleep(time.Second)
		p.SetLEDs([]bool{false, true, false, true, false, true, false, true})
		time.Sleep(time.Second)
	}

}

func testLEDandButtons(p *pkg.LED8KEY) {

	buttons := make([]bool, 8)

	for {

		err := p.ReadButtons(buttons)
		if err != nil {
			fmt.Println("ReadButtons:", err)
		}

		err = p.SetLEDs(buttons)
		if err != nil {
			fmt.Println("SetLEDs:", err)
		}

		time.Sleep(time.Millisecond * 10)

	}

}

func testWriteDigits(p *pkg.LED8KEY) {

	err := p.InitWriteData(true)
	if err != nil {
		fmt.Println("InitWriteData:", err)
	}
	time.Sleep(time.Second)

	data := []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}
	p.WriteDigits(0, data)
}

func testWriteString(p *pkg.LED8KEY) {

	err := p.InitWriteData(true)
	if err != nil {
		fmt.Println("InitWriteData:", err)
	}
	time.Sleep(time.Second)

	p.SetLEDs([]bool{true, true, false, false, true, false, true, false})

	data := "3.1415926"

	err = p.WriteString(0, data)
	if err != nil {
		fmt.Println("WriteString:", err)
	}
}

func main() {

	// Open the rpio once for all using packages (right now just go-led8key)
	fmt.Println("opening gpio")
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	p := pkg.NewLED8KEY(17, 27, 22)

	err = p.ConfigureDisplay(true, 7)
	if err != nil {
		fmt.Println("ConfigureDisplay:", err)
	}

	time.Sleep(time.Second)

	//testLEDs(p)
	//testFill(p)
	//testFlash(p)
	//testHello(p)
	//testButtons(p)
	//testLEDandButtons(p)
	//testWriteDigits(p)
	testWriteString(p)

}

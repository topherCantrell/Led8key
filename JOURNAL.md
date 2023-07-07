# 7/6/2023

I made the constructor function and fixed the package path. 

Now to fetch it:

```
go get github.com/topherCantrell/go-led8key/pkg
```

I'm confused by this. What about a "/internal" directory beside the pkg directory. How do I
install the commands from /cmd. What's all this packaging about? More learning for sure.

Now for some general learning about the chip and its protocol.

I need to separate this into chip driver and a board-specific. The board-specific
references the layouts of things -- LEDs+digits+buttons. The driver is layout agnostic.

# 7/5/2023

Installed GO on the raspberry pi.

```
sudo apt-get install golang

go mod init example.com
go mod tidy

go run led8key.go
```

For now, I'll just use the example.com. Eventually this moves to a real module on my github.

The `led8key.go` started as the sample blinking-LED pin from the tutorial:

https://medium.com/@farissyariati/go-raspberry-pi-hello-world-tutorial-7e830d08b3ae

I verified it worked with an LED on the output pin.

I'd like for the library to work like this:

```go
import github.com/topherCantrell/go-led8key

// board = NewLED8Key(pins)
// board.setLEDs(0b10110101)
```

Much to learn. How does developing a module look along side testing it? Go install?

# 7/10/2023

I added the PrintString with some font characters for numbers and such. The user can
add/change characters in the font. Next I need to tweak the PrintString to combine
decimal points to the character to the left where possible. For instance, "3.14" should
only consume three digits with the decimal point going with the "3". Other things to
handle:
  - ".1234" There is no digit to the left of the "." -- insert a space
  - "2..34" There is no digit to the left of the repeated "." -- insert a space

The combining changes the character count check.

# 7/7/2023

Lots of good progress. Got the reading routine working. I can read the buttons.

I got the display control working. I can flash the intensity on the random garbage
at power up.

I mapped out the LED connections. Time for the "board specific" code

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

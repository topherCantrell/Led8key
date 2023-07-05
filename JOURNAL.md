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
package main

import (
	"github.com/davecheney/junk/rpi/lcd"
)

func main() {
	lcd := lcd.New(
		24, // RS
		25, // enable
		17, 18, 22, 27,
	)
	lcd.Begin(16, 2)
	lcd.Clear()
}

package lcd

import (
	"github.com/davecheney/junk/rpi"
)

type LCD struct {
	rs, enable                                   uint8
	data                                         [4]uint8
	displayfunction, displaycontrol, displaymode uint8
}

const (
	// commands
	LCD_CLEARDISPLAY   = 0x01
	LCD_RETURNHOME     = 0x02
	LCD_ENTRYMODESET   = 0x04
	LCD_DISPLAYCONTROL = 0x08
	LCD_CURSORSHIFT    = 0x10
	LCD_FUNCTIONSET    = 0x20
	LCD_SETCGRAMADDR   = 0x40
	LCD_SETDDRAMADDR   = 0x80

	// flags for display entry mode
	LCD_ENTRYRIGHT          = 0x00
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTINCREMENT = 0x01
	LCD_ENTRYSHIFTDECREMENT = 0x00

	// flags for display on/off control
	LCD_DISPLAYON  = 0x04
	LCD_DISPLAYOFF = 0x00
	LCD_CURSORON   = 0x02
	LCD_CURSOROFF  = 0x00
	LCD_BLINKON    = 0x01
	LCD_BLINKOFF   = 0x00

	// flags for display/cursor shift
	LCD_DISPLAYMOVE = 0x08
	LCD_CURSORMOVE  = 0x00
	LCD_MOVERIGHT   = 0x04
	LCD_MOVELEFT    = 0x00

	// flags for function set
	LCD_8BITMODE = 0x10
	LCD_4BITMODE = 0x00
	LCD_2LINE    = 0x08
	LCD_1LINE    = 0x00
	LCD_5x10DOTS = 0x04
	LCD_5x8DOTS  = 0x00

	HIGH = true
	LOW  = false
)

func New(rs, enable, d0, d1, d2, d3 uint8) *LCD {
	lcd := &LCD{
		rs:              rs,
		enable:          enable,
		data:            [...]uint8{d0, d1, d2, d3},
		displayfunction: LCD_4BITMODE | LCD_1LINE | LCD_5x8DOTS,
		displaymode:     LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT,
	}
	rpi.PinMode(lcd.rs, rpi.OUTPUT)
	rpi.PinMode(lcd.enable, rpi.OUTPUT)
	lcd.Begin(16, 1)
	return lcd
}

func (l *LCD) Begin(cols, lines int) {
	if lines > 1 {
		l.displayfunction |= LCD_2LINE
	}
	rpi.DelayMicroseconds(50000)
	// Now we pull both RS and R/W low to begin commands
	rpi.DigitalWrite(l.rs, LOW)
	rpi.DigitalWrite(l.enable, LOW)
	// we start in 8bit mode, try to set 4 bit mode
	l.write4bits(0x03)
	rpi.DelayMicroseconds(4500) // wait min 4.1ms

	// we start in 8bit mode, try to set 4 bit mode
	l.write4bits(0x03)
	rpi.DelayMicroseconds(4500) // wait min 4.1ms

	// we start in 8bit mode, try to set 4 bit mode
	l.write4bits(0x03)
	rpi.DelayMicroseconds(150)

	// finally, set to 8-bit interface
	l.write4bits(0x02)

	// finally, set # lines, font size, etc.
	l.command(LCD_FUNCTIONSET | l.displayfunction)

	// turn the display on with no cursor or blinking default
	l.displaycontrol = LCD_DISPLAYON | LCD_CURSOROFF | LCD_BLINKOFF
	l.Display()

	// clear it off
	l.Clear()

	// Initialize to default text direction (for romance languages)
	l.displaymode = LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT
	// set the entry mode
	l.command(LCD_ENTRYMODESET | l.displaymode)

}

func (l *LCD) Display() {
	l.displaycontrol |= LCD_DISPLAYON
	l.command(LCD_DISPLAYCONTROL | l.displaycontrol)
}

func (l *LCD) Clear() {
	l.command(LCD_CLEARDISPLAY) // clear display, set cursor position to zero
	rpi.DelayMicroseconds(2000) // this command takes a long time!
}

func (l *LCD) command(v uint8) {
	l.send(v, LOW)
}

func (l *LCD) write(v uint8) {
	l.send(v, HIGH)
}

func (l *LCD) send(v uint8, mode bool) {
	rpi.DigitalWrite(l.rs, mode)
	l.write4bits(v >> 4)
	l.write4bits(v)
}

func (l *LCD) write4bits(v uint8) {
	for i := uint8(0); i < 4; i++ {
		b := (v>>i)&0x01 == 0x01
		rpi.PinMode(l.data[i], rpi.OUTPUT)
		rpi.DigitalWrite(l.data[i], b)
	}
	l.pulseEnable()
}

func (l *LCD) pulseEnable() {
	rpi.DigitalWrite(l.enable, LOW)
	rpi.DelayMicroseconds(1)
	rpi.DigitalWrite(l.enable, HIGH)
	rpi.DelayMicroseconds(1)
	rpi.DigitalWrite(l.enable, LOW)
	rpi.DelayMicroseconds(100)
}

package goeditor

import (
	"github.com/tebeka/atexit"
)

func CtrlKey(k rune) rune {
	return k & 0x1f
}

func MoveCursor(key rune) {
	switch key {
	case 'a':
		config.Cursor.MoveX(-1)
		break
	case 'd':
		config.Cursor.MoveX(1)
		break
	case 'w':
		config.Cursor.MoveY(-1)
		break
	case 's':
		config.Cursor.MoveY(1)
		break
	}
}

func ProcessKeypress() {
	c, err := ReadKey()
	if err != nil {
		atexit.Fatalf("error: %s", err)
		atexit.Exit(1)
	}
	if c == nil {
		atexit.Fatalln("c is nil")
		atexit.Exit(1)
	}
	switch *c {
	case rune('q'):
		RefreshScreen()
		atexit.Exit(0)
		break
	case rune('w'), rune('s'), rune('a'), rune('d'):
		MoveCursor(*c)
	}
}

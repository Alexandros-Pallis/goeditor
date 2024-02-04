package goeditor

import (
	"os"

	"github.com/tebeka/atexit"
)

func CtrlKey(c rune) rune {
	return c & 0x1f
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

func ProcessKeyPress() error {
	char, err := ReadKey()
	if err != nil {
		return err
	}
	switch char {
	case CtrlKey('q'):
		os.Stdout.Write([]byte("\x1b[2J"))
		os.Stdout.Write([]byte("\x1b[H"))
		atexit.Exit(0)
		break
	}
	MoveCursor(char)
	return nil
}

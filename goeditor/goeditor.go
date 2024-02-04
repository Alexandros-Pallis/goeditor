package goeditor

import (
	"github.com/tebeka/atexit"
)

func setup() {
    config.Cursor.X = 0
    config.Cursor.Y = 0
	if err := GetWindowSize(&config.ScreenRows, &config.ScreenCols); err != nil {
		atexit.Fatalf("failed to get window size: %v\r\n", err)
	}
}

func Run() {
	setup()
	if err := EnableRawMode(); err != nil {
		atexit.Fatalf("failed to enable raw mode: %v\r\n", err)
	}
	for {
		RefreshScreen()
		if err := ProcessKeyPress(); err != nil {
			atexit.Fatalf("failed to process key press: %v\r\n", err)
		}
	}
}

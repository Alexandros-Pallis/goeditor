package goeditor

import "github.com/tebeka/atexit"

func setup() {
	config.Cursor.X = 0
	config.Cursor.Y = 0
	if err := GetWindowSize(&config.ScreenRows, &config.ScreenCols); err != nil {
		atexit.Fatalf("failed to get window size: %v\r\n", err)
	}
}
func Run() {
    setup()
	EnableRawMode()
	for {
		RefreshScreen()
		ProcessKeypress()
	}
}

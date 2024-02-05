package goeditor

import (
	"fmt"
	"os"
)

type EditorBuffer struct {
	Bytes []byte
}

func NewBuffer() EditorBuffer {
	return EditorBuffer{
		Bytes: make([]byte, 0),
	}
}
func (buffer *EditorBuffer) Append(b byte) {
	buffer.Bytes = append(buffer.Bytes, b)
}
func (buffer *EditorBuffer) AppendBytes(bytes []byte) {
	for _, b := range bytes {
		buffer.Append(b)
	}
}
func (buffer *EditorBuffer) Reset() {
	buffer.Bytes = make([]byte, 0)
}

var Buffer = NewBuffer()

func DrawRows() {
	for y := 0; y < config.ScreenRows; y++ {
		if y == config.ScreenRows/3 {
			welcome := []byte("Goeditor -- version 1")
			if len(welcome) > config.ScreenCols {
				welcome = welcome[0:config.ScreenCols]
			}
			padding := (config.ScreenCols - len(welcome)) / 2
			for i := 0; i <= padding; i++ {
				if i == 0 {
					Buffer.AppendBytes([]byte("~"))
				} else {
					Buffer.AppendBytes([]byte(" "))
				}
			}
			Buffer.AppendBytes(welcome)
		} else {
			Buffer.AppendBytes([]byte("~"))
		}
		Buffer.AppendBytes([]byte("\x1b[K"))
		if y < config.ScreenRows-1 {
			Buffer.AppendBytes([]byte("\r\n"))
		}
	}
}

func RefreshScreen() {
	Buffer.AppendBytes([]byte("\x1b[?25l"))
	Buffer.AppendBytes([]byte("\x1b[2H"))
	DrawRows()
	cursorPosition := fmt.Sprintf("\x1b[%d;%dH", config.Cursor.Y+1, config.Cursor.X+1)
	Buffer.AppendBytes([]byte(cursorPosition))
	Buffer.AppendBytes([]byte("\x1b[H"))
	Buffer.AppendBytes([]byte("\x1b[?25l"))
	os.Stdout.Write(Buffer.Bytes)
}

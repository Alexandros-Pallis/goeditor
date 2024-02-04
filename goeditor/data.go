package goeditor

import "golang.org/x/sys/unix"

type Cursor struct {
    X int
    Y int
}
func (c *Cursor) MoveX(n int) {
    c.X = c.X + n
}
func (c *Cursor) MoveY(n int) {
    c.Y = c.Y + n
}

type EditorConfig struct {
    OrigTermios unix.Termios
    ScreenRows int
    ScreenCols int
    Cursor Cursor
}

var config EditorConfig


package goeditor

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/pkg/term/termios"
	"github.com/tebeka/atexit"
	"golang.org/x/sys/unix"
)

func EnableRawMode() error {
	termios.Tcgetattr(os.Stdin.Fd(), &config.OrigTermios)
	atexit.Register(disableRawMode)
	config.OrigTermios.Iflag ^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	config.OrigTermios.Oflag ^= unix.OPOST
	config.OrigTermios.Cflag |= unix.CS8
	config.OrigTermios.Lflag ^= unix.ECHO | unix.ICANON | unix.IEXTEN | unix.ISIG
	config.OrigTermios.Cc[unix.VMIN] = 0
	config.OrigTermios.Cc[unix.VTIME] = 1
	err := termios.Tcsetattr(os.Stdin.Fd(), termios.TCSAFLUSH, &config.OrigTermios)
	if err != nil {
		return err
	}
	return nil
}

func disableRawMode() {
	if err := termios.Tcsetattr(os.Stdin.Fd(), termios.TCSAFLUSH, &config.OrigTermios); err != nil {
		log.Fatal(err)
	}
}

func ReadKey() (rune, error) {
	var char rune
	var b = []byte{0}
	n, err := os.Stdin.Read(b)
	if n == 0 || err != nil {
		return char, err
	}
	char = rune(b[0])
	return char, err
}

func clearScreen() {
	os.Stdout.Write([]byte("\x1b[2J"))
	os.Stdout.Write([]byte("\x1b[H"))
}

func GetCursorPosition(rows *int, cols *int) error {
	var buf [32]byte
	number, err := os.Stdout.Write([]byte("\x1b[6n"))
	if number != 4 || err != nil {
		return err
	}
	i := 0
	for i < len(buf)-1 {
		_, err := os.Stdin.Read(buf[i : i+1])
		if err != nil {
			break
		}
		if buf[i] == 'R' {
			break
		}
		i++
	}
	buf[i] = byte(0)
	if buf[0] != '\x1b' || buf[1] != '[' {
		return errors.New("goeditor: buffer doesn't start with expected values")
	}
	if n, err := fmt.Sscanf(string(buf[2]), "%d;%d", rows, cols); n != 2 || err != nil {
		return err
	}
	return nil
}

func GetWindowSize(rows *int, cols *int) error {
	winsize, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil || winsize.Col == 0 {
		fallback := []byte("\x1b[999C\x1b[999B")
		if n, err := os.Stdout.Write(fallback); n != 12 || err != nil {
			return err
		}
		return GetCursorPosition(rows, cols)
	}
	*rows = int(winsize.Row)
	*cols = int(winsize.Col)
	return nil
}

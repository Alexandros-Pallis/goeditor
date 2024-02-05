package goeditor

import (
	"errors"
	"fmt"
	"os"

	"github.com/pkg/term/termios"
	"github.com/tebeka/atexit"
	"golang.org/x/sys/unix"
)

func EnableRawMode() {
	if err := termios.Tcgetattr(os.Stdin.Fd(), &config.OrigTermios); err != nil {
		atexit.Fatalf("tcgetattr: %s", err)
		atexit.Exit(1)
	}
	atexit.Register(disableRawMode)
	config.OrigTermios.Iflag ^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	config.OrigTermios.Oflag ^= unix.OPOST
	config.OrigTermios.Cflag |= unix.CS8
	config.OrigTermios.Lflag ^= unix.ECHO | unix.ICANON | unix.IEXTEN | unix.ISIG
	config.OrigTermios.Cc[unix.VMIN] = 0
	config.OrigTermios.Cc[unix.VTIME] = 1
	if err := termios.Tcsetattr(os.Stdin.Fd(), termios.TCSAFLUSH, &config.OrigTermios); err != nil {
		atexit.Fatalf("tcsetattr: %s", err)
		atexit.Exit(1)
	}
}

func disableRawMode() {
	if err := termios.Tcsetattr(os.Stdin.Fd(), termios.TCSAFLUSH, &config.OrigTermios); err != nil {
		atexit.Fatalf("tcsetattr: %s", err)
		atexit.Exit(1)
	}
}

func ReadKey() (*rune, error) {
	bytes := make([]byte, 1)
	var c rune
	for {
		nread, err := os.Stdin.Read(bytes)
		if err != nil && err.Error() != "EOF" {
			atexit.Fatalf("error reading bytes: %s", err)
			atexit.Exit(1)
		}
		if nread == 1 {
			c = rune(bytes[0])
			break
		}
	}
	return &c, nil
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

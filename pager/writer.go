package pager

import (
	"io"
	"os"
)

// SetW
type SetW interface {
	SetStdW(io.Writer)
	SetErrW(io.Writer)
	StdW() io.Writer
	ErrW() io.Writer
}

// Writers provides a canonical implementation of the SetW interface. It
// can be used as a mixin with other data structures.
type Writers struct {
	stdW io.Writer
	errW io.Writer
}

// W creates and returns a Writers object initialising it with a
// the os.Stdout and os.Stderr values
func W() Writers {
	return Writers{
		stdW: os.Stdout,
		errW: os.Stderr,
	}
}

// SetStdW sets the value of the standard writer
func (ws *Writers) SetStdW(w io.Writer) { ws.stdW = w }

// SetErrW sets the value of the error writer
func (ws *Writers) SetErrW(w io.Writer) { ws.errW = w }

// StdW returns the value of the standard writer
func (ws Writers) StdW() io.Writer { return ws.stdW }

// ErrW returns the value of the error writer
func (ws Writers) ErrW() io.Writer { return ws.errW }

package pager

import (
	"io"
	"os"
)

// SetWriter
type SetWriter interface {
	SetStdWriter(io.Writer)
	SetErrWriter(io.Writer)
	StdWriter() io.Writer
	ErrWriter() io.Writer
}

// Writers provides a canonical implementation of the SetWriter interface. It
// can be used as a mixin with other data structures.
type Writers struct {
	stdW io.Writer
	errW io.Writer
}

// DfltWriters creates and returns a Writers object initialising it with a
// the os.Stdout and os.Stderr values
func DfltWriters() Writers {
	return Writers{
		stdW: os.Stdout,
		errW: os.Stderr,
	}
}

// SetStdWriter sets the value of the standard writer
func (ws *Writers) SetStdWriter(w io.Writer) { ws.stdW = w }

// SetErrWriter sets the value of the error writer
func (ws *Writers) SetErrWriter(w io.Writer) { ws.errW = w }

// StdWriter returns the value of the standard writer
func (ws Writers) StdWriter() io.Writer { return ws.stdW }

// ErrWriter returns the value of the error writer
func (ws Writers) ErrWriter() io.Writer { return ws.errW }

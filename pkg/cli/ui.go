package cli

import (
	"fmt"
	"io"
)

// UI is an interface for interacting with the terminal, or "interface"
// of a CLI. This abstraction doesn't have to be used, but helps provide
// a simple, layerable way to manage user interactions.
type UI interface {
	// Input uses the query to ask users for input. The response is
	// returned as the given string, or an error.
	// Input(string) (string, error)

	// Output is called for normal standard output.
	Output(string)

	// Info is called for information related to the previous output.
	// In general this may be the exact same as Output, but this gives
	// Ui implementors some flexibility with output formats.
	Info(string)

	// Error is used for any error messages that might appear on standard
	// error.
	Error(string)

	// Warn is used for any warning messages that might appear on standard
	// error.
	Warn(string)
}

// BasicUI is an implementation of UI that outputs to the given writer.
type SimpleUI struct {
	Reader      io.Reader
	Writer      io.Writer
	ErrorWriter io.Writer
}

func (u *SimpleUI) Error(message string) {
	w := u.Writer
	if u.ErrorWriter != nil {
		w = u.ErrorWriter
	}
	fmt.Fprint(w, message, "\n")
}

func (u *SimpleUI) Info(message string) {
	u.Output(message)
}

func (u *SimpleUI) Output(message string) {
	fmt.Fprint(u.Writer, message, "\n")
}

func (u *SimpleUI) Warn(message string) {
	u.Error(message)
}

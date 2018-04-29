package imagemagick

import "fmt"

// ParserError represents an error by the parser
type ParserError struct {
	msg    string
	file   string
	cmd    string
	stdOut []byte
	stdErr []byte
}

// NewParserError creates a new ParserError
func NewParserError(msg string, file string, cmd string, stdOut []byte, stdErr []byte) *ParserError {
	return &ParserError{
		msg:    msg,
		file:   file,
		cmd:    cmd,
		stdOut: stdOut,
		stdErr: stdErr,
	}
}

// Error returns a string representation of the error with all of its properties
func (err *ParserError) Error() string {
	if err == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Error: %v; File: %q; Cmd: %q; StdOut: %q; StdErr: %q",
		err.msg,
		err.file,
		err.cmd,
		err.stdOut,
		err.stdErr,
	)
}

// Msg returns the error message
func (err *ParserError) Msg() string {
	return err.msg
}

// File returns the file that caused the error (if any)
func (err *ParserError) File() string {
	return err.file
}

// Cmd returns the command that caused the error (if any)
func (err *ParserError) Cmd() string {
	return err.cmd
}

// StdOut that was produced by the failed command (if any)
func (err *ParserError) StdOut() []byte {
	return err.stdOut
}

// StdErr that was produced by the failed command (if any)
func (err *ParserError) StdErr() []byte {
	return err.stdErr
}

package imagemagick

import "fmt"

type ParserError struct {
	msg    string
	file   string
	cmd    string
	stdOut []byte
	stdErr []byte
}

func NewParserError(msg string, file string, cmd string, stdOut []byte, stdErr []byte) *ParserError {
	return &ParserError{
		msg:    msg,
		file:   file,
		cmd:    cmd,
		stdOut: stdOut,
		stdErr: stdErr,
	}
}

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

func (err *ParserError) Msg() string {
	return err.msg
}

func (err *ParserError) File() string {
	return err.file
}

func (err *ParserError) Cmd() string {
	return err.cmd
}

func (err *ParserError) StdOut() []byte {
	return err.stdOut
}

func (err *ParserError) StdErr() []byte {
	return err.stdErr
}

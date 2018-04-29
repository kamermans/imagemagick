package imagemagick

import "fmt"

type ImageMagickParserError struct {
	msg    string
	file   string
	cmd    string
	stdOut []byte
	stdErr []byte
}

func NewImageMagickParserError(msg string, file string, cmd string, stdOut []byte, stdErr []byte) *ImageMagickParserError {
	return &ImageMagickParserError{
		msg:    msg,
		file:   file,
		cmd:    cmd,
		stdOut: stdOut,
		stdErr: stdErr,
	}
}

func (err *ImageMagickParserError) Error() string {
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

func (err *ImageMagickParserError) Msg() string {
	return err.msg
}

func (err *ImageMagickParserError) File() string {
	return err.file
}

func (err *ImageMagickParserError) Cmd() string {
	return err.cmd
}

func (err *ImageMagickParserError) StdOut() []byte {
	return err.stdOut
}

func (err *ImageMagickParserError) StdErr() []byte {
	return err.stdErr
}

package imagemagick_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kamermans/imagemagick"
)

func getTestError() *imagemagick.ParserError {
	msg := "test message"
	file := "fname"
	cmd := "./foo bar"
	stdOut := []byte("foo works!")
	stdErr := []byte("foo failed!")

	return imagemagick.NewParserError(msg, file, cmd, stdOut, stdErr)
}

func TestNewParserError(t *testing.T) {

	err := getTestError()
	if err == nil {
		t.Fatalf("Got nil error")
	}

}

func TestNewParserErrorActsLikeAnError(t *testing.T) {

	testFunc := func(e error) string {
		return fmt.Sprintf("%v", e)
	}

	err := getTestError()
	result := testFunc(err)

	expectedResult := `StdErr: "foo failed!"`
	if !strings.Contains(result, expectedResult) {
		t.Fatalf("Error message did not contain expected string %q, got %q", expectedResult, result)
	}
}

func TestNewParserErrorAvoidsNilPointerr(t *testing.T) {

	var err *imagemagick.ParserError

	expectedMsg := "<nil>"
	if err.Error() != expectedMsg {
		t.Fatalf("Nil error returned unexpected message: expected %v, got %v", expectedMsg, err.Error())
	}
}

func TestParserErrorGetters(t *testing.T) {

	err := getTestError()

	msg := "test message"
	file := "fname"
	cmd := "./foo bar"
	stdOut := []byte("foo works!")
	stdErr := []byte("foo failed!")

	if err.Msg() != msg {
		t.Fatalf("Msg() getter failed")
	}

	if err.File() != file {
		t.Fatalf("File() getter failed")
	}

	if err.Cmd() != cmd {
		t.Fatalf("Cmd() getter failed")
	}

	if string(err.StdOut()) != string(stdOut) {
		t.Fatalf("StdOut() getter failed")
	}

	if string(err.StdErr()) != string(stdErr) {
		t.Fatalf("StdErr() getter failed")
	}

}

// Package imagemagick provides a high-level wrapper for the ImageMagick
// `convert` command and a replacement for the `identify` command to gather
// detail information on images like width, height, exif tags, colorspace, etc.
package imagemagick

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// Parser represents an ImageMagick command-line tool parser
type Parser struct {
	// The 'convert' command
	ConvertCommand string
	// Number of files to pass to convert at once when running in parallel
	BatchSize int
	// Number of workers to start when running in parallel (default: # of CPUs)
	Workers int

	// Used to clean the ImageMagick JSON
	jsonCleaner     *regexp.Regexp
	jsonCleanerRepl []byte

	// Used for testing
	command func(name string, arg ...string) *exec.Cmd
}

// Used to clean the dirty ImageMagick JSON
const (
	jsonCleanerPattern    = ": -?(?:1.#IN[DF]|nan|inf)"
	jsonCleanerReplString = ": null"
)

// NewParser creates a new Parser object
func NewParser() *Parser {
	return &Parser{
		ConvertCommand:  "convert",
		BatchSize:       20,
		Workers:         runtime.NumCPU(),
		jsonCleaner:     regexp.MustCompile(jsonCleanerPattern),
		jsonCleanerRepl: []byte(jsonCleanerReplString),
		command:         exec.Command,
	}
}

// SetCommand allows you to set an alternate exec.Cmd object, which is useful for mocking
// commands for testing
func (parser *Parser) SetCommand(command func(name string, arg ...string) *exec.Cmd) {
	parser.command = command
}

// GetImageDetailsParallel computes ImageDetails for a channel of input files.  The results
// are available in the results channel and errors are on the errors channel.  You should read
// the results and errors channels in a go routine to prevent blocking.  The number of workers
// is defined at Parser.Workers.  ImageMagick supports batches of input files, and this function
// uses batches of size Parse.BatchSize.  When a batch of files is passed to ImageMagick and an
// error is encountered, the batch is split up and each file is sent individually so the bad
// file can be identified and sent to the errors channel.
func (parser *Parser) GetImageDetailsParallel(
	files <-chan string,
	results chan<- *ImageResult,
	errs chan<- *ParserError,
) {
	go func() {
		sendImageDetails := func(fileBatch ...string) {
			detailsSlice, err := parser.GetImageDetails(fileBatch...)

			if err != nil {
				if len(fileBatch) == 1 {
					errs <- err
					return
				}

				// Reprocess this batch one-by-one since at least one of the files failed
				// and caused the whole batch to be lost
				detailsSlice = []*ImageResult{}
				for _, file := range fileBatch {
					thisFileDetails, thisErr := parser.GetImageDetails(file)
					if thisErr != nil {
						errs <- thisErr
						continue
					}

					detailsSlice = append(detailsSlice, thisFileDetails...)
				}
			}

			for _, details := range detailsSlice {
				results <- details
			}
		}

		var wg sync.WaitGroup
		wg.Add(parser.Workers)
		for w := 0; w < parser.Workers; w++ {

			go func() {
				// Collect a batch of files to pass to ImageMagick
				fileBatch := []string{}

				for file := range files {
					fileBatch = append(fileBatch, file)
					if len(fileBatch) == parser.BatchSize {
						sendImageDetails(fileBatch...)
						fileBatch = []string{}
					}
				}

				if len(fileBatch) > 0 {
					sendImageDetails(fileBatch...)
				}

				wg.Done()
			}()

		}

		// Wait for all the workers to finish, then close the results channel
		wg.Wait()
		close(results)
	}()
}

// GetImageDetails computes ImageDetails for one or more input files, returning (results, err).
// If an error is encountered, results will be nil and err will contain the error.
func (parser *Parser) GetImageDetails(files ...string) (results []*ImageResult, err *ParserError) {
	// Compose command like this:
	//   "convert file1 file2 fileN json:-"
	args := append(files, "json:-")
	cmd := parser.command(parser.ConvertCommand, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if cmdErr := cmd.Run(); cmdErr != nil {
		cmdParts := []string{parser.ConvertCommand}
		cmdParts = append(cmdParts, files...)
		cmdParts = append(cmdParts, "json:-")

		err = NewParserError(
			"ImageMagick convert command failed: "+cmdErr.Error(),
			strings.Join(files, ", "),
			strings.Join(cmdParts, " "),
			stdout.Bytes(),
			stderr.Bytes(),
		)
		return
	}

	jsonBlob := stdout.Bytes()

	results, jsonErr := parser.GetImageDetailsFromJSON(&jsonBlob)
	if jsonErr != nil {
		err = NewParserError(
			jsonErr.Error(),
			strings.Join(files, ", "),
			"",
			[]byte{},
			[]byte{},
		)
	}

	return
}

// GetImageDetailsFromJSON computes ImageDetails for the given JSON data, returning (results, err).
// If an error is encountered, results will be nil and err will contain the error.  Note that the
// JSON data is cleaned of invalid numbers with Regexp because ImageMagick `convert` leaks C++ NaNs
// into the output data, like `{"bytes": -nan}` and `{"entropy": -1.#IND}`
func (parser *Parser) GetImageDetailsFromJSON(jsonBlob *[]byte) (results []*ImageResult, err error) {

	// Clean up C++ NaNs on Windows. On Linux/Unix, C++ produces nan and inf, which get parsed correctly
	// ImageMagick leaks this non-stanard JSON in fields like channelStatistics.Alpha.standardDeviation
	jsonBlobObj := parser.jsonCleaner.ReplaceAll(*jsonBlob, parser.jsonCleanerRepl)
	jsonBlob = &jsonBlobObj

	//	jsonBlob, _ := ioutil.ReadFile(fname)
	err = json.Unmarshal(*jsonBlob, &results)
	if err != nil {
		err = fmt.Errorf("Unable to decode ImageMagick JSON: %v", err)
	}
	return
}

// Convert is a helper to call the ImageMagick `convert` command.  It will return the stdOut, stdErr and
// a ParserError if the command failed (by returing a non-zero exit code, for example)
func (parser *Parser) Convert(args ...string) (stdOut *[]byte, stdErr *[]byte, err *ParserError) {

	cmd := parser.command(parser.ConvertCommand, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if cmdErr := cmd.Run(); cmdErr != nil {
		cmdParts := []string{parser.ConvertCommand}
		cmdParts = append(cmdParts, args...)

		err = NewParserError(
			"ImageMagick convert command failed: "+cmdErr.Error(),
			"",
			strings.Join(cmdParts, " "),
			stdout.Bytes(),
			stderr.Bytes(),
		)
	}

	stdOutBytes := stdout.Bytes()
	stdErrBytes := stderr.Bytes()

	return &stdOutBytes, &stdErrBytes, err

}

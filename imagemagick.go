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

type ImageMagickParser struct {
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

const (
	jsonCleanerPattern    = ": -?(?:1.#IN[DF]|nan|inf)"
	jsonCleanerReplString = ": null"
)

func NewImageMagickParser() *ImageMagickParser {
	return &ImageMagickParser{
		ConvertCommand:  "convert",
		BatchSize:       20,
		Workers:         runtime.NumCPU(),
		jsonCleaner:     regexp.MustCompile(jsonCleanerPattern),
		jsonCleanerRepl: []byte(jsonCleanerReplString),
		command:         exec.Command,
	}
}

func (parser *ImageMagickParser) SetCommand(command func(name string, arg ...string) *exec.Cmd) {
	parser.command = command
}

func (parser *ImageMagickParser) GetImageDetailsParallel(
	files <-chan string,
	results chan<- *ImageMagickDetails,
	errs chan<- *ImageMagickParserError,
) (done chan bool) {

	done = make(chan bool)

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
				detailsSlice = []*ImageMagickDetails{}
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

	return done
}

func (parser *ImageMagickParser) GetImageDetails(files ...string) (results []*ImageMagickDetails, err *ImageMagickParserError) {
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

		err = NewImageMagickParserError(
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
		err = NewImageMagickParserError(
			jsonErr.Error(),
			strings.Join(files, ", "),
			"",
			[]byte{},
			[]byte{},
		)
	}

	return
}

func (parser *ImageMagickParser) GetImageDetailsFromJSON(jsonBlob *[]byte) (results []*ImageMagickDetails, err error) {

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

func (parser *ImageMagickParser) Convert(args ...string) (stdOut *[]byte, stdErr *[]byte, err *ImageMagickParserError) {

	cmd := parser.command(parser.ConvertCommand, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if cmdErr := cmd.Run(); cmdErr != nil {
		cmdParts := []string{parser.ConvertCommand}
		cmdParts = append(cmdParts, args...)

		err = NewImageMagickParserError(
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

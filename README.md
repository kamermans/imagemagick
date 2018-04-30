# kamermans/imagemagick

[![Build Status](https://travis-ci.org/kamermans/imagemagick.svg?branch=master)](https://travis-ci.org/kamermans/imagemagick) [![Go Report Card](https://goreportcard.com/badge/github.com/kamermans/imagemagick)](https://goreportcard.com/report/github.com/kamermans/imagemagick) [![Godoc](https://img.shields.io/badge/godoc-docs-blue.svg)](https://godoc.org/github.com/kamermans/imagemagick)

High-level Go wrapper for the ImageMagick `convert` command and a replacement for the `identify` command to gather detailed information on images like width, height, exif tags, colorspace, etc, without requiring the ImageMagick shared libraries; works on Linux, Windows, MacOS and any other system that has access to the `convert` command.

Check the godocs for the latest documentation:
https://godoc.org/github.com/kamermans/imagemagick

## Usage
For usage examples, check out the [godocs](https://godoc.org/github.com/kamermans/imagemagick) and the `examples_*.go` files.

One of the great features of `kamermans/imagemagick` is that it comes with built-in support for processing files in parallel:

```go
package main

package imagemagick_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kamermans/imagemagick"
)

func main() {
    var (
        convertCmd    = `/usr/local/bin/convert`
        imageFilesDir = `/tmp/sample_images`
    )

    parser := imagemagick.NewParser()
    parser.ConvertCommand = convertCmd

    files := make(chan string)
    results := make(chan *imagemagick.ImageResult)
    errs := make(chan *imagemagick.ParserError)

    // Used to tell us when the results have all be consumed
    done := make(chan bool)

    parser.GetImageDetailsParallel(files, results, errs)

    // Send in files
    go func() {
        defer close(files)

        filepath.Walk(imageFilesDir, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                fmt.Printf("Unable to access path %q: %v\n", imageFilesDir, err)
                return err
            }

            if info.IsDir() {
                return nil
            }

            // Send this image into the files channel
            files <- path

            return nil
        })
    }()

    numErrors := 0
    numResults := 0
    startTime := time.Now()

    // Store the number of images of each format that we've seen
    resultsByFormat := map[string]int64{}
    // Store the total size of the images we've seen
    totalSize := int64(0)
    // Report progress this often
    reportInterval := 2 * time.Second

    // Report progress
    go func() {
        time.Sleep(reportInterval)
        for {
            // Get a sorted list of formats so it looks consistent
            formats := []string{}
            for format := range resultsByFormat {
                formats = append(formats, format)
            }
            sort.Strings(formats)

            outLines := []string{}
            for _, format := range formats {
                outLines = append(outLines, fmt.Sprintf("%v: %v", format, resultsByFormat[format]))
            }

            numPerSecond := float64(numResults+numErrors) / time.Since(startTime).Seconds()
            fmt.Printf("Results: %v, Errors: %v, Rate: %.0f/sec, Image Data: %v MB, Formats: {%v}\n",
                numResults,
                numErrors,
                numPerSecond,
                totalSize/1000000,
                strings.Join(outLines, ", "),
            )

            time.Sleep(reportInterval)
        }
    }()

    // Consume results and errors
    moreErrs := true
    moreResults := true
    for {
        if !moreErrs && !moreResults {
            break
        }

        select {
        case _, ok := <-errs:
            if !ok {
                moreErrs = false
                continue
            }
            numErrors++
        case details, ok := <-results:
            if !ok {
                moreResults = false
                continue
            }
            numResults++
            image := details.Image

            // Collect some stats for the progress function above
            totalSize += image.Size()
            if details.Image.Format != "" {
                resultsByFormat[details.Image.Format]++
            }

            // You can get the image details here if you want
            // fmt.Printf("Received result for image: %v (%v)\n",
            // 	image.BaseName,
            // 	image.Format,
            // )
        }
    }

    fmt.Printf("Received %v results and %v errors\n", numResults, numErrors)

    // Here's what the output looks like on my laptop with 4523 sample images:
    //
    // Results: 40, Errors: 0, Rate: 20/sec, Image Data: 3 MB, Formats: {JPEG: 40}
    // Results: 160, Errors: 0, Rate: 40/sec, Image Data: 9 MB, Formats: {JPEG: 159, PNG: 1}
    // Results: 280, Errors: 0, Rate: 46/sec, Image Data: 18 MB, Formats: {JPEG: 279, PNG: 1}
    // ... lots of output ...
    // Results: 4304, Errors: 16, Rate: 46/sec, Image Data: 303 MB, Formats: {JPEG: 4281, PNG: 23}
    // Results: 4386, Errors: 17, Rate: 46/sec, Image Data: 305 MB, Formats: {JPEG: 4362, PNG: 24}
    // Results: 4465, Errors: 19, Rate: 46/sec, Image Data: 309 MB, Formats: {JPEG: 4440, PNG: 25}
    // Received 4503 results and 20 errors
}
```

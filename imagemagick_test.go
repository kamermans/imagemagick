package imagemagick_test

import (
    "github.com/kamermans/imagemagick/test_resources"
)

func TestGetImageDetailsFromJSON(t *testing.T) {
    files, err := filepath.Glob("test_resources/json_output/*.json")
    if err != nil || len(files) == 0 {
        t.Fatalf("Cannot read JSON test files")
    }

    parser := imagemagick.NewImageMagickParser()

    for _, file := range files {
        shouldFail := strings.Contains(file, "_error_")
        jsonBlob, readErr := ioutil.ReadFile(file)
        if readErr != nil {
            t.Fatalf("Cannot read JSON test file %q: %v", file, readErr.Error())
        }

        results, jsonErr := parser.GetImageDetailsFromJSON(&jsonBlob)

        if shouldFail {
            if jsonErr == nil {
                t.Fatalf("JSON unmarshal was expected to fail but it did not")
            }
        } else {
            if jsonErr != nil {
                t.Fatalf("Cannot decode JSON from test file %q: %v", file, jsonErr.Error())
            }

            if len(results) == 0 {
                t.Fatalf("Cannot decode JSON from test file %q: nil error was returned with no results", file)
            }

            if strings.Contains(path.Base(file), "multi") && len(results) < 2 {
                t.Fatalf("Cannot decode JSON from test file %q: multiple input resulted in one output", file)
            }
        }
    }
}

func TestGetImageDetailsProperties(t *testing.T) {

    file := "test_resources/json_output/image_metadata_exif_linux.json"
    jsonBlob, readErr := ioutil.ReadFile(file)
    if readErr != nil {
        t.Fatalf("Cannot read JSON test file %q: %v", file, readErr.Error())
    }

    parser := imagemagick.NewImageMagickParser()

    results, jsonErr := parser.GetImageDetailsFromJSON(&jsonBlob)
    if jsonErr != nil {
        t.Fatalf("Cannot decode JSON from test file %q: %v", file, jsonErr.Error())
    }

    if len(results) == 0 {
        t.Fatalf("Cannot decode JSON from test file %q: nil error was returned with no results", file)
    }

    image := results[0].Image
    props := image.PropertiesMap()
    exProp := 5
    if len(props) != exProp {
        t.Fatalf("Image properties are wrong length: expected %v, got %v", exProp, len(props))
    }

    tagTypes := []string{}
    for tagType := range props {
        tagTypes = append(tagTypes, tagType)
    }
    sort.Strings(tagTypes)

    exTagTypes := []string{"date", "exif", "icc", "jpeg", "signature"}
    for i := range tagTypes {
        if exTagTypes[i] != tagTypes[i] {
            t.Fatalf("Image properties tagTypes are wrong: expected %v, got %v", exTagTypes, tagTypes)
        }
    }

    exif := image.ExifTags()

    exSoftware := "Adobe Photoshop CC 2017 (Macintosh)"
    acSoftware, ok := exif["Software"]
    if !ok || acSoftware != exSoftware {
        t.Fatalf("Image exif data is wrong or missing: expected %v, got %v", exSoftware, acSoftware)
    }
}

func TestGetImageDetails(t *testing.T) {
    file := "test_resources/json_output/image_metadata_multi_formats_linux2.json"
    mockExec := test_resources.NewMockExec("TestHelperGetImageDetails")

    parser := imagemagick.NewImageMagickParser()
    parser.SetCommand(mockExec.Command)

    d, err := parser.GetImageDetails(file)
    if err != nil {
        t.Fatalf("GetImageDetails() failed: %v", err.Error())
    }

    if mockExec.RunCount() != 1 {
        t.Fatalf("GetImageDetails() failed: command was not run")
    }

    lastRun := mockExec.LastRun()

    expectedCmd := "convert"
    actualCmd := lastRun.GetCommand()
    if expectedCmd != actualCmd {
        t.Fatalf("GetImageDetails() failed: expected %v, got %v", expectedCmd, actualCmd)
    }

    expectedArgs := []string{file, "json:-"}
    actualArgs := lastRun.GetArgs()
    if len(expectedArgs) != len(actualArgs) {
        t.Fatalf("GetImageDetails() failed: expected %v, got %v", expectedArgs, actualArgs)
    }

    for i := range expectedArgs {
        if expectedArgs[i] != actualArgs[i] {
            t.Fatalf("GetImageDetails() failed: expected %v, got %v", expectedArgs, actualArgs)
        }
    }

    if len(d) != 40 {
        t.Fatalf("GetImageDetails() failed")
    }

}

func TestGetImageDetailsCustomConvertCommand(t *testing.T) {
    file := "test_resources/json_output/image_metadata_multi_formats_linux2.json"
    mockExec := test_resources.NewMockExec("TestHelperGetImageDetails")

    parser := imagemagick.NewImageMagickParser()
    parser.SetCommand(mockExec.Command)
    parser.ConvertCommand = "foobar.exe"

    _, err := parser.GetImageDetails(file)
    if err != nil {
        t.Fatalf("GetImageDetails() failed: %v", err.Error())
    }

    if mockExec.RunCount() != 1 {
        t.Fatalf("GetImageDetails() failed: command was not run")
    }

    lastRun := mockExec.LastRun()

    expectedCmd := "foobar.exe"
    actualCmd := lastRun.GetCommand()
    if expectedCmd != actualCmd {
        t.Fatalf("GetImageDetails() failed: expected %v, got %v", expectedCmd, actualCmd)
    }
}

func TestConvert(t *testing.T) {
    mockExec := test_resources.NewMockExec("TestHelperGetImageDetails")

    expectedCmd := "foobar.exe"
    expectedArgs := []string{"foo", "--bar", "baz", "name:Bilbo Baggins"}

    parser := imagemagick.NewImageMagickParser()
    parser.SetCommand(mockExec.Command)
    parser.ConvertCommand = expectedCmd

    actualStdOut, actualStdErr, err := parser.Convert(expectedArgs...)
    if err != nil {
        t.Fatalf("Convert() failed: %v", err.Error())
    }

    file := "test_resources/json_output/image_metadata_multi_formats_linux2.json"
    expectedStdOut, readErr := ioutil.ReadFile(file)
    if readErr != nil {
        t.Fatalf("Can't read test resource file %v: %v", file, readErr.Error())
    }

    if len(*actualStdErr) != 0 {
        t.Fatalf("Convert() failed: expected StdErr to be empty, got: %v", string(*actualStdErr))
    }

    if string(*actualStdOut) != string(expectedStdOut) {
        t.Fatalf("Convert() failed: StdOut does not match, expected: %v got: %v", string(expectedStdOut), string(*actualStdOut))
    }

    if mockExec.RunCount() != 1 {
        t.Fatalf("Convert() failed: command was not run")
    }

    lastRun := mockExec.LastRun()
    actualArgs := lastRun.GetArgs()
    if len(expectedArgs) != len(actualArgs) {
        t.Fatalf("Convert() failed: expected %v, got %v", expectedArgs, actualArgs)
    }

    for i := range expectedArgs {
        if expectedArgs[i] != actualArgs[i] {
            t.Fatalf("Convert() failed: expected %v, got %v", expectedArgs, actualArgs)
        }
    }

    actualCmd := lastRun.GetCommand()
    if expectedCmd != actualCmd {
        t.Fatalf("Convert() failed: expected %v, got %v", expectedCmd, actualCmd)
    }
}

func TestHelperGetImageDetails(t *testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
        return
    }
    defer os.Exit(0)

    file := "test_resources/json_output/image_metadata_multi_formats_linux2.json"
    jsonBlob, readErr := ioutil.ReadFile(file)
    if readErr != nil {
        fmt.Fprintf(os.Stderr, "%s\n", readErr.Error())
        os.Exit(2)
    }

    os.Stdout.Write(jsonBlob)

}

func TestGetImageDetailsFailed(t *testing.T) {
    mockExec := test_resources.NewMockExec("TestHelperGetImageDetailsFailed")
    parser := imagemagick.NewImageMagickParser()
    parser.SetCommand(mockExec.Command)

    file := "test_resources/json_output/image_metadata_multi_formats_linux2.json"

    _, err := parser.GetImageDetails(file)
    if err == nil {
        t.Fatalf("GetImageDetails() did not fail as expected")
    }

    expectedStdErr := "Simulated Failure"
    if !strings.Contains(err.Error(), expectedStdErr) {
        t.Fatalf("GetImageDetails() failed: expected %v, got %v", expectedStdErr, err.Error())
    }

    if mockExec.RunCount() != 1 {
        t.Fatalf("GetImageDetails() failed: command was not run")
    }

}

func TestConvertFailed(t *testing.T) {
    mockExec := test_resources.NewMockExec("TestHelperGetImageDetailsFailed")

    expectedCmd := "foobar.exe"
    expectedArgs := []string{"foo", "--bar", "baz", "name:Bilbo Baggins"}

    parser := imagemagick.NewImageMagickParser()
    parser.SetCommand(mockExec.Command)
    parser.ConvertCommand = expectedCmd

    _, actualStdErr, err := parser.Convert(expectedArgs...)
    if err == nil {
        t.Fatalf("Convert() did not fail as expected")
    }

    if mockExec.RunCount() != 1 {
        t.Fatalf("Convert() failed: command was not run")
    }

    lastRun := mockExec.LastRun()
    actualArgs := lastRun.GetArgs()
    if len(expectedArgs) != len(actualArgs) {
        t.Fatalf("Convert() failed: expected %v, got %v", expectedArgs, actualArgs)
    }

    for i := range expectedArgs {
        if expectedArgs[i] != actualArgs[i] {
            t.Fatalf("Convert() failed: expected %v, got %v", expectedArgs, actualArgs)
        }
    }

    actualCmd := lastRun.GetCommand()
    if expectedCmd != actualCmd {
        t.Fatalf("Convert() failed: expected %v, got %v", expectedCmd, actualCmd)
    }

    expectedStdErr := "Simulated Failure"
    if string(*actualStdErr) != expectedStdErr {
        t.Fatalf("Convert() failed: expected %v, got %v", expectedStdErr, string(*actualStdErr))
    }
}

func TestHelperGetImageDetailsFailed(t *testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
        return
    }
    defer os.Exit(2)

    file := "test_resources/json_output/image_metadata_multi_formats_linux2.json"
    jsonBlob, readErr := ioutil.ReadFile(file)
    if readErr != nil {
        fmt.Fprintf(os.Stderr, "%s\n", readErr.Error())
    }

    os.Stdout.Write(jsonBlob)
    os.Stderr.Write([]byte("Simulated Failure"))
}

func TestGetImageDetailsParallel(t *testing.T) {
    mockExec := test_resources.NewMockExec("TestHelperGetImageDetailsParallel")

    parser := imagemagick.NewImageMagickParser()
    parser.SetCommand(mockExec.Command)
    parser.Workers = 8
    parser.BatchSize = 4

    files := make(chan string)
    results := make(chan *imagemagick.ImageMagickDetails)
    errs := make(chan *imagemagick.ImageMagickParserError)

    done := parser.GetImageDetailsParallel(files, results, errs)

    const numTestFiles = 40
    testFiles := [numTestFiles]string{}
    for i := 0; i < numTestFiles; i++ {
        testFiles[i] = fmt.Sprintf("/foo/bar/test_file/%d", i)
    }

    // Send in files
    go func() {
        for _, testFile := range testFiles {
            files <- testFile
    	}
    	close(files)
    }()

    // Consume errors
    receivedErrors := []*imagemagick.ImageMagickParserError{}
    go func() {
    	for err := range errs {
            receivedErrors = append(receivedErrors, err)
    	}
    }()

    // Read out results
    receivedImages := []*imagemagick.ImageMagickImageDetails{}
    go func() {
    	for details := range results {
            receivedImages = append(receivedImages, details.Image)
    	}
    	done <- true
    	close(errs)
    }()

    <-done

    expectedRuns := numTestFiles / parser.BatchSize
    actualRuns := mockExec.RunCount()
    if expectedRuns != actualRuns {
        t.Fatalf("GetImageDetailsParallel() failed: wrong run count: expected %v got %v", expectedRuns, actualRuns)
    }

    expectedCmd := "convert"

    testFilesMap := map[string]bool{}
    for _, testFileName := range testFiles {
        testFilesMap[testFileName] = false
    }

    // Trigger Duplicate
    // testFilesMap["/foo/bar/test_file/7"] = true

    // Trigger Not Received
    // testFilesMap["/foo/bar/missing_file/0"] = false

    // Trigger Received Unrequested
    // delete(testFilesMap, "/foo/bar/test_file/7")

    for _, run := range mockExec.Runs() {
        actualCmd := run.GetCommand()
        if expectedCmd != actualCmd {
            t.Fatalf("GetImageDetailsParallel() failed: expected %v, got %v", expectedCmd, actualCmd)
        }

        args := run.GetArgs()
        lastArg := len(args) - 1

        for a := 0; a < lastArg; a++ {

            file := args[a]

            if _, ok := testFilesMap[file]; !ok {
                t.Fatalf("GetImageDetailsParallel() failed: received unrequested file: %v", file)
                continue
            }

            if testFilesMap[file] {
                t.Fatalf("GetImageDetailsParallel() failed: received duplicate file: %v", file)
                continue
            }

            testFilesMap[file] = true
        }

        expectedArg := "json:-"
        if args[lastArg] != expectedArg {
            t.Fatalf("GetImageDetailsParallel() failed: expected %v, got %v", expectedArg, args[lastArg])
        }
    }

    for f, ok := range testFilesMap {
        if !ok {
            t.Fatalf("GetImageDetailsParallel() failed: file was not received: %v", f)
        }
    }

    if len(receivedErrors) != 0 {
        t.Fatalf("GetImageDetailsParallel() failed: received %v errors", len(receivedErrors))
    }

    if len(receivedImages) != numTestFiles {
        t.Fatalf("GetImageDetailsParallel() failed: expected %v results, got %v", numTestFiles, len(receivedImages))
    }

}

func TestHelperGetImageDetailsParallel(t *testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
        return
    }
    defer os.Exit(0)

    // There are 4 results in this file
    file := "test_resources/json_output/image_metadata_multi_linux.json"
    jsonBlob, readErr := ioutil.ReadFile(file)
    if readErr != nil {
        fmt.Fprintf(os.Stderr, "%s\n", readErr.Error())
        os.Exit(2)
    }

    os.Stdout.Write(jsonBlob)

}

func TestGetImageDetailsParallelWithErrors(t *testing.T) {
    mockExec := test_resources.NewMockExec("TestHelperGetImageDetailsParallelWithErrors")

    parser := imagemagick.NewImageMagickParser()
    parser.SetCommand(mockExec.Command)
    parser.Workers = 8
    parser.BatchSize = 4

    files := make(chan string)
    results := make(chan *imagemagick.ImageMagickDetails)
    errs := make(chan *imagemagick.ImageMagickParserError)

    done := parser.GetImageDetailsParallel(files, results, errs)

    const numTestFiles = 40

    testFiles := [numTestFiles]string{}
    for i := 0; i < numTestFiles; i++ {
        testFiles[i] = fmt.Sprintf("/foo/bar/test_file/%d", i)
    }

    // Send in files
    go func() {
        for _, testFile := range testFiles {
            files <- testFile
    	}
    	close(files)
    }()

    // Consume errors
    receivedErrors := []*imagemagick.ImageMagickParserError{}
    go func() {
    	for err := range errs {
            receivedErrors = append(receivedErrors, err)
    	}
    }()

    // Read out results
    receivedImages := []*imagemagick.ImageMagickImageDetails{}
    go func() {
    	for details := range results {
            receivedImages = append(receivedImages, details.Image)
    	}
    	done <- true
    	close(errs)
    }()

    <-done

    actualRuns := mockExec.RunCount()
    if actualRuns < 40 {
        t.Fatalf("GetImageDetailsParallel() failed: wrong run count: expected >= 40 got %v", actualRuns)
    }

    expectedCmd := "convert"

    testFilesMap := map[string]int{}
    for _, testFileName := range testFiles {
        testFilesMap[testFileName] = 0
    }

    // Trigger Duplicate
    // testFilesMap["/foo/bar/test_file/7"] = true

    // Trigger Not Received
    // testFilesMap["/foo/bar/missing_file/0"] = false

    // Trigger Received Unrequested
    // delete(testFilesMap, "/foo/bar/test_file/7")

    for _, run := range mockExec.Runs() {
        actualCmd := run.GetCommand()
        if expectedCmd != actualCmd {
            t.Fatalf("GetImageDetailsParallel() failed: expected %v, got %v", expectedCmd, actualCmd)
        }

        args := run.GetArgs()
        lastArg := len(args) - 1

        for a := 0; a < lastArg; a++ {

            file := args[a]

            if _, ok := testFilesMap[file]; !ok {
                t.Fatalf("GetImageDetailsParallel() failed: received unrequested file: %v", file)
                continue
            }

            testFilesMap[file]++
        }

        expectedArg := "json:-"
        if args[lastArg] != expectedArg {
            t.Fatalf("GetImageDetailsParallel() failed: expected %v, got %v", expectedArg, args[lastArg])
        }
    }

    for f, count := range testFilesMap {
        if count == 0 {
            t.Fatalf("GetImageDetailsParallel() failed: file was not received: %v", f)
        }
    }

    if len(receivedErrors) < 30 {
        t.Fatalf("GetImageDetailsParallel() failed: expected >= 30 errors, got %v", len(receivedErrors))
    }

    // Note that receivedImages is 0 here because the retried images all failed too
    // due to the mocking system, but we cen't carry state into the mock because it's
    // done in a completely different process.

}

func TestHelperGetImageDetailsParallelWithErrors(t *testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
        return
    }
    defer os.Exit(0)

    // There are 4 results in this file and 1 of them will cause a parsing error
    file := "test_resources/json_output/image_metadata_multi_error_linux.json"
    jsonBlob, readErr := ioutil.ReadFile(file)
    if readErr != nil {
        fmt.Fprintf(os.Stderr, "%s\n", readErr.Error())
        os.Exit(2)
    }

    os.Stdout.Write(jsonBlob)

}

package imagemagick_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kamermans/imagemagick"
)

// An example of getting image details from an image file.  For this example we're downloading the
// file from w3.org, writing it to a temp file and getting its details
func Example_GetImageDetails() {

	var (
		convertCmd = `/usr/loca/bin/convert`
		imageURL   = `https://www.w3.org/People/mimasa/test/imgformat/img/w3c_home.gif`
	)

	imageFile, err := downloadTempImage(imageURL)
	if err != nil {
		panic("Could not download the example image")
	}
	defer os.Remove(imageFile)

	parser := imagemagick.NewParser()
	parser.ConvertCommand = convertCmd
	results, detErr := parser.GetImageDetails(imageFile)
	if detErr != nil {
		panic(detErr.Error())
	}

	// Note that one output JSON can contain multiple results
	image := results[0].Image

	// Print the format
	fmt.Printf("Format: %v (%v)\n", image.Format, image.MimeType)

	// Print the geometry
	fmt.Printf("Dimensions: %v\n", *image.Geometry.Dimensions)

	// Example output:
	// Format: GIF (image/gif)
	// Dimensions: {Width: 72, Height: 48}

}

func downloadTempImage(imageURL string) (file string, err error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fh, err := ioutil.TempFile("", "imagemagick_example_")
	if err != nil {
		return
	}
	defer fh.Close()

	_, err = fh.Write(body)
	if err != nil {
		return
	}

	file = fh.Name()

	return
}

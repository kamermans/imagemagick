package imagemagick_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kamermans/imagemagick"
)

func TestImageResult(t *testing.T) {
	d := &imagemagick.ImageResult{
		Image: new(imagemagick.ImageDetails),
	}

	if d.Image == nil {
		t.Fatalf("Image should not be nil")
	}
}

func TestImageDetails(t *testing.T) {
	d := new(imagemagick.ImageDetails)

	if d.Name != "" {
		t.Fatalf("Name should be empty")
	}
}

func TestDetailsToJSON(t *testing.T) {
	d := &imagemagick.ImageResult{
		Image: new(imagemagick.ImageDetails),
	}

	jBytes, err := d.ToJSON(false)
	if err != nil {
		t.Fatalf("JSON conversion failed: %v", err.Error())
	}

	if !bytes.HasPrefix(jBytes, []byte(`{"image":{`)) {
		t.Fatalf("JSON output is invalid: %v", string(jBytes))
	}

	if bytes.Contains(jBytes, []byte("\n")) {
		t.Fatalf("Non-pretty JSON should not contain a newline")
	}
}

func TestDetailsToJSONPretty(t *testing.T) {
	d := &imagemagick.ImageResult{
		Image: new(imagemagick.ImageDetails),
	}

	jBytes, err := d.ToJSON(true)
	if err != nil {
		t.Fatalf("JSON conversion failed: %v", err.Error())
	}

	if !bytes.Contains(jBytes, []byte(`"image": {`)) {
		t.Fatalf("JSON output is invalid: %v", string(jBytes))
	}

	if !bytes.Contains(jBytes, []byte("\n")) {
		t.Fatalf("Pretty JSON should contain newlines")
	}
}

func TestImageDetailsToJSON(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Name: "testimage.jpg",
		Geometry: &imagemagick.Geometry{
			Dimensions: &imagemagick.Dimensions{
				Width:  320,
				Height: 240,
			},
		},
	}

	jBytes, err := d.ToJSON(false)
	if err != nil {
		t.Fatalf("JSON conversion failed: %v", err.Error())
	}

	if !bytes.Contains(jBytes, []byte(`"width":320`)) {
		t.Fatalf("JSON output is invalid: %v", string(jBytes))
	}

	if bytes.Contains(jBytes, []byte("\n")) {
		t.Fatalf("Non-pretty JSON should not contain a newline")
	}
}

func TestImageDetailsToJSONPretty(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Name: "testimage.jpg",
		Geometry: &imagemagick.Geometry{
			Dimensions: &imagemagick.Dimensions{
				Width:  320,
				Height: 240,
			},
		},
	}

	jBytes, err := d.ToJSON(true)
	if err != nil {
		t.Fatalf("JSON conversion failed: %v", err.Error())
	}

	if !bytes.Contains(jBytes, []byte(`"width": 320`)) {
		t.Fatalf("JSON output is invalid: %v", string(jBytes))
	}

	if !bytes.Contains(jBytes, []byte("\n")) {
		t.Fatalf("Pretty JSON should contain newlines")
	}
}

func TestImageDetailsCanvaseDimensions(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Name: "testimage.jpg",
		Geometry: &imagemagick.Geometry{
			Point: &imagemagick.Point{
				X: 100,
				Y: 100,
			},
			Dimensions: &imagemagick.Dimensions{
				Width:  320,
				Height: 240,
			},
		},
	}

	e1 := int64(420)
	if d.Geometry.CanvasWidth() != e1 {
		t.Fatalf("CanvasWidth() failed, expected %v, got %v", e1, d.Geometry.CanvasWidth())
	}

	e2 := int64(340)
	if d.Geometry.CanvasHeight() != e2 {
		t.Fatalf("CanvasHeight() failed, expected %v, got %v", e2, d.Geometry.CanvasHeight())
	}

}

func TestImageDetailsGeoOffset(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Name: "testimage.jpg",
		Geometry: &imagemagick.Geometry{
			Point: &imagemagick.Point{
				X: 100,
				Y: 100,
			},
			Dimensions: &imagemagick.Dimensions{
				Width:  320,
				Height: 240,
			},
		},
	}

	e := &imagemagick.Point{
		X: 100,
		Y: 100,
	}

	a := d.Geometry.Offset()

	if a == nil {
		t.Fatalf("Offset() returned nil")
	}

	if e.X != a.X || e.Y != a.Y {
		t.Fatalf("Offset() failed, expected %v, got %v", e, a)
	}
}

func TestImageDetailsSize(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Filesize: "1200B",
	}

	e := int64(1200)
	a := d.Size()
	if a != e {
		t.Fatalf("Size() failed, expected %v, got %v", e, a)
	}
}

func TestProfileTotalSize(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Filesize: "1200B",
		Profiles: map[string]map[string]interface{}{
			"foo": {
				"length": int64(256),
			},
			"bar": {
				"length": int(512),
			},
			"baz": {
				"length": 0,
			},
			"floatyval": {
				"length": 128.75,
			},
		},
	}

	e := int64(896)
	a := d.ProfileTotalSize()
	if a != e {
		t.Fatalf("ProfileTotalSize() failed, expected %v, got %v", e, a)
	}
}

func TestProfileSizePercent(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Filesize: "1200B",
		Profiles: map[string]map[string]interface{}{
			"foo": {
				"length": int64(256),
			},
			"bar": {
				"length": int(512),
			},
			"baz": {
				"length": 0,
			},
			"floatyval": {
				"length": 128.75,
			},
		},
	}

	e := "0.4275"
	a := fmt.Sprintf("%.4f", d.ProfileSizePercent())
	if a != e {
		t.Fatalf("ProfileSizePercent() failed, expected %v, got %v", e, a)
	}
}

func TestHasProfile(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Filesize: "1200B",
		Profiles: map[string]map[string]interface{}{
			"foo": {
				"length": int64(256),
			},
			"bar": {
				"length": int(512),
			},
			"baz": {
				"length": 0,
			},
			"floatyval": {
				"length": 128.75,
			},
		},
	}

	checks := map[string]bool{
		"foo":       true,
		"something": false,
		"baz":       false, // because it's empty
	}

	for name, e := range checks {
		a := d.HasProfile(name)
		if a != e {
			t.Fatalf("HasProfile() failed for %v, expected %v, got %v", name, e, a)
		}
	}
}

func TestProfileNames(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Filesize: "1200B",
		Profiles: map[string]map[string]interface{}{
			"foo": {
				"length": int64(256),
			},
			"bar": {
				"length": int(512),
			},
			"baz": {
				"length": 0,
			},
			"floatyval": {
				"length": 128.75,
			},
		},
	}

	e := []string{"bar", "baz", "floatyval", "foo"}
	a := d.ProfileNames()

	if len(e) != len(a) {
		t.Fatalf("ProfileNames() failed, expected %v, got %v", e, a)
	}

	for i, v := range a {
		if v != e[i] {
			t.Fatalf("ProfileNames() failed, expected %v, got %v", e, a)
		}
	}
}

func TestProfileSizes(t *testing.T) {
	d := &imagemagick.ImageDetails{
		Filesize: "1200B",
		Profiles: map[string]map[string]interface{}{
			"foo": {
				"length": int64(256),
			},
			"bar": {
				"length": int(512),
			},
			"baz": {
				"length": 0,
			},
			"floatyval": {
				"length": 128.75,
			},
		},
	}

	e := map[string]int64{
		"foo":       int64(256),
		"bar":       int64(512),
		"baz":       int64(0),
		"floatyval": int64(128),
	}
	a := d.ProfileSizes()

	if len(e) != len(a) {
		t.Fatalf("ProfileSizes() failed, expected %v, got %v", e, a)
	}

	for i, v := range a {
		if v != e[i] {
			t.Fatalf("ProfileSizes() failed, expected %v, got %v", e, a)
		}
	}
}

func TestGeometry(t *testing.T) {
	d := &imagemagick.Geometry{
		&imagemagick.Point{
			X: 15,
			Y: 20,
		},
		&imagemagick.Dimensions{
			Width:  147,
			Height: 239,
		},
	}

	e := int64(162)
	a := d.CanvasWidth()
	if e != a {
		t.Fatalf("CanvasWidth() failed, expected %v, got %v", e, a)
	}

	e = int64(259)
	a = d.CanvasHeight()
	if e != a {
		t.Fatalf("CanvasHeight() failed, expected %v, got %v", e, a)
	}

	aPoint := d.Offset()
	ePoint := imagemagick.Point{
		X: 15,
		Y: 20,
	}
	if ePoint.X != aPoint.X || ePoint.Y != aPoint.Y {
		t.Fatalf("Offset() failed, expected %v, got %v", ePoint, aPoint)
	}
}

//

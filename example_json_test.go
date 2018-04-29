package imagemagick_test

import (
	"fmt"
	"sort"

	"github.com/kamermans/imagemagick"
)

// Example provides an example of getting image details from the ImageMagick
// `convert` tool's `json` format.
func Example() {

	jsonBlob := getSampleJSONOutput()

	parser := imagemagick.NewParser()
	results, err := parser.GetImageDetailsFromJSON(jsonBlob)
	if err != nil {
		panic(err.Error())
	}

	// Note that one output JSON can contain multiple results
	image := results[0].Image

	// Print the filename
	fmt.Printf("Image: %v (%v)\n", image.BaseName, image.Format)

	// Collect the ICC information
	props := image.PropertiesMap()
	icc := props["icc"]
	// Sort ICC tags and print them in order so the tests don't fail :)
	propTags := []string{}
	for tag := range icc {
		propTags = append(propTags, tag)
	}
	sort.Strings(propTags)

	// Print the ICC information
	for _, tag := range propTags {
		fmt.Printf("ICC %v: %v\n", tag, icc[tag])
	}

	// Print the geometry
	fmt.Printf("Geometry: %v\n", *image.Geometry)

	// Print the profile size
	profileSizePct := image.ProfileSizePercent() * 100.0
	fmt.Printf("Profiles: The embedded profiles account for %.2f%% of the total file size\n", profileSizePct)

	// Output:
	// Image: bug72278.jpg (JPEG)
	// ICC copyright: Copyright (c) 1998 Hewlett-Packard Company
	// ICC description: sRGB IEC61966-2.1
	// ICC manufacturer: IEC http://www.iec.ch
	// ICC model: IEC 61966-2.1 Default RGB colour space - sRGB
	// Geometry: {{X: 0, Y: 0} {Width: 300, Height: 300}}
	// Profiles: The embedded profiles account for 16.48% of the total file size
}

func getSampleJSONOutput() *[]byte {
	data := []byte(sampleJSONOutput)
	return &data
}

// This is an example of the output you get from ImageMagick `convert`
// when you use `json` for the output format.  For example, to get the
// JSON details for `foo.jpg`, you would use this command:
//      convert foo.jpg foo_details.json
// This package uses the STDOUT method to avoid writing an output file:
//      convert foo.jpg json:-
const sampleJSONOutput = `[{
"image": {
 "name": "json:/tmp/image_metadata_multi_formats_linux.json",
 "baseName": "bug72278.jpg",
 "format": "JPEG",
 "formatDescription": "JPEG",
 "mimeType": "image/jpeg",
 "class": "DirectClass",
 "geometry": {
    "width": 300,
    "height": 300,
    "x": 0,
    "y": 0
 },
 "resolution": {
    "x": 200,
    "y": 200
 },
 "printSize": {
    "x": 1.5,
    "y": 1.5
 },
 "units": "PixelsPerInch",
 "type": "Bilevel",
 "baseType": "Undefined",
 "endianess": "Undefined",
 "colorspace": "sRGB",
 "depth": 1,
 "baseDepth": 8,
 "channelDepth": {
    "red": 1,
    "green": 1,
    "blue": 1
 },
 "pixels": 270000,
 "imageStatistics": {
    "Overall": {
      "min": 255,
      "max": 255,
      "mean": 255,
      "standardDeviation": 0,
      "kurtosis": 1.6384e+64,
      "skewness": 9.375e+44,
      "entropy": -nan
    }
 },
 "channelStatistics": {
    "Red": {
      "min": 255,
      "max": 255,
      "mean": 255,
      "standardDeviation": 0,
      "kurtosis": 8.192e+63,
      "skewness": 1e+45,
      "entropy": -nan
    },
    "Green": {
      "min": 255,
      "max": 255,
      "mean": 255,
      "standardDeviation": 0,
      "kurtosis": 8.192e+63,
      "skewness": 1e+45,
      "entropy": -nan
    },
    "Blue": {
      "min": 255,
      "max": 255,
      "mean": 255,
      "standardDeviation": 0,
      "kurtosis": 8.192e+63,
      "skewness": 1e+45,
      "entropy": -nan
    }
 },
 "renderingIntent": "Perceptual",
 "gamma": 0.454545,
 "chromaticity": {
    "redPrimary": {
      "x": 0.64,
      "y": 0.33
    },
    "greenPrimary": {
      "x": 0.3,
      "y": 0.6
    },
    "bluePrimary": {
      "x": 0.15,
      "y": 0.06
    },
    "whitePrimary": {
      "x": 0.3127,
      "y": 0.329
    }
 },
 "matteColor": "#BDBDBD",
 "backgroundColor": "#FFFFFF",
 "borderColor": "#DFDFDF",
 "transparentColor": "#00000000",
 "interlace": "None",
 "intensity": "Undefined",
 "compose": "Over",
 "pageGeometry": {
    "width": 300,
    "height": 300,
    "x": 0,
    "y": 0
 },
 "dispose": "Undefined",
 "iterations": 0,
 "scene": 13,
 "scenes": 26,
 "compression": "None",
 "quality": 79,
 "orientation": "Undefined",
 "properties": {
    "comment": "Test",
    "date:create": "2017-10-19T10:30:02-04:00",
    "date:modify": "2017-10-19T10:30:02-04:00",
    "exif:ColorSpace": "1",
    "exif:ComponentsConfiguration": "1, 2, 3, 0",
    "exif:Copyright": "Test",
    "exif:DateTime": "2008:04:03 11:06:23",
    "exif:ExifImageLength": "300",
    "exif:ExifImageWidth": "300",
    "exif:ExifOffset": "196",
    "exif:ExifVersion": "48, 50, 50, 48",
    "exif:FlashPixVersion": "48, 49, 48, 48",
    "exif:ResolutionUnit": "2",
    "exif:Software": "Paint Shop Pro Photo 12.00",
    "exif:thumbnail:Compression": "6",
    "exif:thumbnail:JPEGInterchangeFormat": "380",
    "exif:thumbnail:JPEGInterchangeFormatLength": "1325",
    "exif:thumbnail:ResolutionUnit": "2",
    "exif:thumbnail:XResolution": "787399/10000",
    "exif:thumbnail:YCbCrPositioning": "2",
    "exif:thumbnail:YResolution": "787399/10000",
    "exif:XResolution": "1999995/10000",
    "exif:YCbCrPositioning": "2",
    "exif:YResolution": "1999995/10000",
    "icc:copyright": "Copyright (c) 1998 Hewlett-Packard Company",
    "icc:description": "sRGB IEC61966-2.1",
    "icc:manufacturer": "IEC http://www.iec.ch",
    "icc:model": "IEC 61966-2.1 Default RGB colour space - sRGB",
    "jpeg:colorspace": "2",
    "jpeg:sampling-factor": "1x1,1x1,1x1",
    "signature": "31fed455c2bb6e7258a946a2adc33d8493c2084346b9bb000e5042977c56221e"
 },
 "profiles": {
    "8bim": {
      "length": 28
    },
    "exif": {
      "length": 1717
    },
    "icc": {
      "length": 7261
    },
    "iptc": {
      "Unknown[2,0]": [null],
      "Copyright String[2,116]": ["Test"],
      "length": 16
    }
 },
 "tainted": false,
 "filesize": "45720B",
 "numberPixels": "90000",
 "pixelsPerSecond": "692308B",
 "userTime": "0.150u",
 "elapsedTime": "0:01.129",
 "version": "/usr/local/share/doc/ImageMagick-7//index.html"
}
}]`

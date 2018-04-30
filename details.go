package imagemagick

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ImageResult is the top-Level result from ImageMagick.  You almost certainly want to access the Image
// property, but this wrapper is left here for future use, should other types by introduced
type ImageResult struct {
	Image *ImageDetails `json:"image"`
}

// ImageDetails provides detailed information on the image, there are many helpful methods on this object
type ImageDetails struct {
	Alpha             string                            `json:"alpha"`               //"#00FF0000"
	BackgroundColor   string                            `json:"backgroundColor"`     //"#FFFFFF"
	BaseDepth         int64                             `json:"baseDepth"`           //8
	BaseName          string                            `json:"baseName"`            //"image_0002c93a9c0c53e7379a4524fa953ebb"
	BaseType          string                            `json:"baseType"`            //"Undefined"
	BorderColor       string                            `json:"borderColor"`         //"#DFDFDF"
	ChannelDepth      map[string]int64                  `json:"channelDepth"`        //
	ChannelStatistics map[string]*ChannelStatistics     `json:"channelStatistics"`   //
	Chromaticity      map[string]*PointFloat            `json:"chromaticity"`        //
	Class             string                            `json:"class"`               //"DirectClass"
	Colormap          []string                          `json:"colormap"`            //["#7F82B8FF","#393747FF"]
	ColormapEntries   int64                             `json:"colormapEntries"`     //128
	Colorspace        string                            `json:"colorspace"`          //"sRGB"
	Compose           string                            `json:"compose"`             //"Over"
	Compression       string                            `json:"compression"`         //"JPEG2000"
	Depth             int64                             `json:"depth"`               //8
	Dispose           string                            `json:"dispose"`             //"Undefined"
	ElapsedTime       string                            `json:"elapsedTime"`         //"0:01.049"
	Endianess         string                            `json:"endianess"`           //"Undefined"
	Filesize          string                            `json:"filesize"`            //"0B"
	Format            string                            `json:"format"`              //"JP2"
	FormatDescription string                            `json:"formatDescription"`   //"JP2"
	Gamma             float64                           `json:"gamma"`               //0.454545
	Geometry          *Geometry                         `json:"geometry"`            //
	ImageStatistics   map[string]*ChannelStatistics     `json:"imageStatistics"`     //
	Intensity         string                            `json:"intensity"`           //"Undefined"
	Interlace         string                            `json:"interlace"`           //"None"
	Iterations        int64                             `json:"iterations"`          //0
	MatteColor        string                            `json:"matteColor"`          //"#BDBDBD"
	MimeType          string                            `json:"mimeType"`            //"image/jp2"
	Name              string                            `json:"name"`                //"test.json"
	NumberPixels      int64                             `json:"numberPixels,string"` //"211750"
	Orientation       string                            `json:"orientation"`         //"Undefined"
	PageGeometry      *Geometry                         `json:"pageGeometry"`        //
	Pixels            int64                             `json:"pixels"`              //635250
	PixelsPerSecond   string                            `json:"pixelsPerSecond"`     //"4235000B"
	PrintSize         *PointFloat                       `json:"printSize"`           //{"x": 2.08333,"y": 1.04167}
	Profiles          map[string]map[string]interface{} `json:"profiles"`            //
	Properties        map[string]string                 `json:"properties"`          //
	Quality           int64                             `json:"quality"`             //75
	RenderingIntent   string                            `json:"renderingIntent"`     //"Perceptual"
	Resolution        *PointFloat                       `json:"resolution"`          //{"x": 96,"y": 96}
	Scene             int64                             `json:"scene"`               //12
	Scenes            int64                             `json:"scenes"`              //26
	Tainted           bool                              `json:"tainted"`             //false
	TransparentColor  string                            `json:"transparentColor"`    //"#00000000"
	Type              string                            `json:"type"`                //"TrueColor"
	Units             string                            `json:"units"`               //"Undefined"
	UserTime          string                            `json:"userTime"`            //"0.030u"
	Version           string                            `json:"version"`             //"/usr/local/share/doc/ImageMagick-7//index.html"
}

// Size of the image in bytes. ImageMagick returns a strangely-formatted string and this the in64 equivalent
func (details ImageDetails) Size() int64 {
	size, _ := strconv.Atoi(strings.Trim(details.Filesize, "B"))
	return int64(size)
}

// ProfileSizePercent returns the percentage of the total filesize which is used by the profiles
func (details ImageDetails) ProfileSizePercent() float64 {
	return float64(details.ProfileTotalSize()) / float64(details.ProfileTotalSize()+details.Size())
}

// ProfileNames of the embedded profiles.  Note that all profiles are included, even
// if they are zero-length
func (details ImageDetails) ProfileNames() (names []string) {
	for name := range details.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return
}

// HasProfile returns true if the image has an embedded profile of the given type.
// Possible options include, but are not limited to: 8bim, exif, iptc, xmp, icc, app1, app12
// Note that zero-length profiles will return false
func (details ImageDetails) HasProfile(name string) bool {
	sizes := details.ProfileSizes()
	if size, ok := sizes[name]; ok {
		return size > 0
	}
	return false
}

// ExifTags returns a map of EXIF tags to their values.  These are pulled from the Properties slice.  Note that
// the prefix "exif:" is timmed from the tag name.  An empty map is returned if there are no
// EXIF tags present
func (details ImageDetails) ExifTags() map[string]string {
	data := details.PropertiesMap("exif")
	if exif, ok := data["exif"]; ok {
		return exif
	}
	return map[string]string{}
}

// PropertiesMap returns a map of the image Properties. The key is split on the first ":" and grouped by the
// first half (the tag name) so the map is a map of map[string]string like this:
//	{
//		"icc": {
//			"brand": "Canon",
//			"model": "EOS 5D Mark IV",
//		},
//		"exif": {
//			"Software": "Adobe Photoshop CC 2017 (Macintosh)",
//		},
//	}
func (details ImageDetails) PropertiesMap(tagFilter ...string) map[string]map[string]string {

	// Easy lookup map to see if filter is in the list
	tagFilterLookup := map[string]bool{}
	for _, filter := range tagFilter {
		tagFilterLookup[filter] = true
	}

	props := map[string]map[string]string{}
	for name, value := range details.Properties {
		parts := strings.SplitN(name, ":", 2)
		tagType := parts[0]

		// If tag filtering is enabled and this tag isn't in the filter, skip it
		if len(tagFilter) > 0 && !tagFilterLookup[tagType] {
			continue
		}

		if _, ok := props[tagType]; !ok {
			props[tagType] = map[string]string{}
		}

		if len(parts) == 1 {
			continue
		}

		tag := parts[1]

		props[tagType][tag] = value
	}

	return props
}

// ProfileSizes returns a map of embedded profile names to their size in bytes
func (details ImageDetails) ProfileSizes() (lengths map[string]int64) {
	lengths = map[string]int64{}
	for name, props := range details.Profiles {
		if rawVal, ok := props["length"]; ok {
			switch val := rawVal.(type) {
			case int:
				lengths[name] = int64(val)
			case int64:
				lengths[name] = val
			case float64:
				lengths[name] = int64(val)
			}
		}
	}
	return
}

// ProfileTotalSize returns the total byte size of all the embedded profiles
func (details ImageDetails) ProfileTotalSize() (size int64) {
	for _, s := range details.ProfileSizes() {
		size += s
	}
	return
}

// Point represents an X, Y point / coordinate
type Point struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

// String representation
func (p Point) String() string {
	return fmt.Sprintf("{X: %v, Y: %v}", p.X, p.Y)
}

// PointFloat represents a float64 X, Y point / coordinate
type PointFloat struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// String representation
func (p PointFloat) String() string {
	return fmt.Sprintf("{X: %v, Y: %v}", p.X, p.Y)
}

// Dimensions represents box dimensions with Width and Height
type Dimensions struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

// String representation
func (d Dimensions) String() string {
	return fmt.Sprintf("{Width: %v, Height: %v}", d.Width, d.Height)
}

// Geometry represents image geometry, including a Point{X, Y} offset and Dimensions{Width, Height} dimensions
type Geometry struct {
	*Point
	*Dimensions
}

// Canvas is the total width and height of the canvas (width/height + x/y offset)
func (geo Geometry) Canvas() *Dimensions {
	return &Dimensions{
		Width:  geo.Width + geo.X,
		Height: geo.Height + geo.Y,
	}
}

// Offset coordinates of the box on the canvas
func (geo Geometry) Offset() *Point {
	return geo.Point
}

// ChannelStatistics represents the image color channel statistics
type ChannelStatistics struct {
	Min               float64 `json:"min"`               // 0,
	Max               float64 `json:"max"`               // 255,
	Mean              float64 `json:"mean"`              // 187.475,
	StandardDeviation float64 `json:"standardDeviation"` // 90.9415,
	Kurtosis          float64 `json:"kurtosis"`          // -1.22588,
	Skewness          float64 `json:"skewness"`          // -0.755169,
	Entropy           float64 `json:"entropy"`           // 0.515529
}

// ToJSON returns the JSON representation of this object
func (details *ImageResult) ToJSON(pretty bool) (out []byte, err error) {
	if pretty {
		return json.MarshalIndent(details, "", "  ")
	}
	return json.Marshal(details)
}

// ToJSON returns the JSON representation of this object
func (details *ImageDetails) ToJSON(pretty bool) (out []byte, err error) {
	if pretty {
		return json.MarshalIndent(details, "", "  ")
	}
	return json.Marshal(details)
}

package imagemagick

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ImageMagickDetails struct {
	Image *ImageMagickImageDetails `json:"image"`
}

type ImageMagickImageDetails struct {
	Alpha             string                                   `json:"alpha"`               //"#00FF0000"
	BackgroundColor   string                                   `json:"backgroundColor"`     //"#FFFFFF"
	BaseDepth         int64                                    `json:"baseDepth"`           //8
	BaseName          string                                   `json:"baseName"`            //"image_0002c93a9c0c53e7379a4524fa953ebb"
	BaseType          string                                   `json:"baseType"`            //"Undefined"
	BorderColor       string                                   `json:"borderColor"`         //"#DFDFDF"
	ChannelDepth      map[string]int64                         `json:"channelDepth"`        //
	ChannelStatistics map[string]*ImageMagickChannelStatistics `json:"channelStatistics"`   //
	Chromaticity      map[string]*ImageMagickPointFloat        `json:"chromaticity"`        //
	Class             string                                   `json:"class"`               //"DirectClass"
	Colormap          []string                                 `json:"colormap"`            //["#7F82B8FF","#393747FF"]
	ColormapEntries   int64                                    `json:"colormapEntries"`     //128
	Colorspace        string                                   `json:"colorspace"`          //"sRGB"
	Compose           string                                   `json:"compose"`             //"Over"
	Compression       string                                   `json:"compression"`         //"JPEG2000"
	Depth             int64                                    `json:"depth"`               //8
	Dispose           string                                   `json:"dispose"`             //"Undefined"
	ElapsedTime       string                                   `json:"elapsedTime"`         //"0:01.049"
	Endianess         string                                   `json:"endianess"`           //"Undefined"
	Filesize          string                                   `json:"filesize"`            //"0B"
	Format            string                                   `json:"format"`              //"JP2"
	FormatDescription string                                   `json:"formatDescription"`   //"JP2"
	Gamma             float64                                  `json:"gamma"`               //0.454545
	Geometry          *ImageMagickGeometry                     `json:"geometry"`            //
	ImageStatistics   map[string]*ImageMagickChannelStatistics `json:"imageStatistics"`     //
	Intensity         string                                   `json:"intensity"`           //"Undefined"
	Interlace         string                                   `json:"interlace"`           //"None"
	Iterations        int64                                    `json:"iterations"`          //0
	MatteColor        string                                   `json:"matteColor"`          //"#BDBDBD"
	MimeType          string                                   `json:"mimeType"`            //"image/jp2"
	Name              string                                   `json:"name"`                //"test.json"
	NumberPixels      int64                                    `json:"numberPixels,string"` //"211750"
	Orientation       string                                   `json:"orientation"`         //"Undefined"
	PageGeometry      *ImageMagickGeometry                     `json:"pageGeometry"`        //
	Pixels            int64                                    `json:"pixels"`              //635250
	PixelsPerSecond   string                                   `json:"pixelsPerSecond"`     //"4235000B"
	PrintSize         *ImageMagickPointFloat                   `json:"printSize"`           //{"x": 2.08333,"y": 1.04167}
	Profiles          map[string]map[string]interface{}        `json:"profiles"`            //
	Properties        map[string]string                        `json:"properties"`          //
	Quality           int64                                    `json:"quality"`             //75
	RenderingIntent   string                                   `json:"renderingIntent"`     //"Perceptual"
	Resolution        *ImageMagickPointFloat                   `json:"resolution"`          //{"x": 96,"y": 96}
	Scene             int64                                    `json:"scene"`               //12
	Scenes            int64                                    `json:"scenes"`              //26
	Tainted           bool                                     `json:"tainted"`             //false
	TransparentColor  string                                   `json:"transparentColor"`    //"#00000000"
	Type              string                                   `json:"type"`                //"TrueColor"
	Units             string                                   `json:"units"`               //"Undefined"
	UserTime          string                                   `json:"userTime"`            //"0.030u"
	Version           string                                   `json:"version"`             //"/usr/local/share/doc/ImageMagick-7//index.html"
}

// Filesize - this is the int64 version of the strangely-formatted strings that
// ImageMagick returns
func (details ImageMagickImageDetails) Size() int64 {
	size, _ := strconv.Atoi(strings.Trim(details.Filesize, "B"))
	return int64(size)
}

// The percentage of the total filesize which is used by the profiles
func (details ImageMagickImageDetails) ProfileSizePercent() float64 {
	return float64(details.ProfileTotalSize()) / float64(details.ProfileTotalSize()+details.Size())
}

// Returns a slice of the embedded profiles.  Note that all profiles are included, even
// if they are zero-length
func (details ImageMagickImageDetails) ProfileNames() (names []string) {
	for name := range details.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return
}

// Returns true if the image has an embedded profile of the given type.
// Possible options include, but are not limited to: 8bim, exif, iptc, xmp, icc, app1, app12
// Note that zero-length profiles will return false
func (details ImageMagickImageDetails) HasProfile(name string) bool {
	sizes := details.ProfileSizes()
	if size, ok := sizes[name]; ok {
		return size > 0
	}
	return false
}

func (details ImageMagickImageDetails) ExifTags() map[string]string {
	data := details.PropertiesMap("exif")
	if exif, ok := data["exif"]; ok {
		return exif
	}
	return map[string]string{}
}

func (details ImageMagickImageDetails) PropertiesMap(tagFilter ...string) map[string]map[string]string {

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
		if len(tagFilter) > 0 {
			if _, ok := tagFilterLookup[tagType]; !ok {
				continue
			}
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

// Returns a map of the embedded profile name to its size in bytes
func (details ImageMagickImageDetails) ProfileSizes() (lengths map[string]int64) {
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

// Returns the total byte size of all the embedded profiles
func (details ImageMagickImageDetails) ProfileTotalSize() (size int64) {
	for _, s := range details.ProfileSizes() {
		size += s
	}
	return
}

type ImageMagickPoint struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

func (p ImageMagickPoint) String() string {
	return fmt.Sprintf("{X: %v, Y: %v}", p.X, p.Y)
}

type ImageMagickPointFloat struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (p ImageMagickPointFloat) String() string {
	return fmt.Sprintf("{X: %v, Y: %v}", p.X, p.Y)
}

type ImageMagickDimensions struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

func (d ImageMagickDimensions) String() string {
	return fmt.Sprintf("{Width: %v, Height: %v}", d.Width, d.Height)
}

type ImageMagickGeometry struct {
	*ImageMagickPoint
	*ImageMagickDimensions
}

func (geo ImageMagickGeometry) CanvasWidth() int64 {
	return geo.Width + geo.X
}

func (geo ImageMagickGeometry) CanvasHeight() int64 {
	return geo.Height + geo.Y
}

func (geo ImageMagickGeometry) Offset() *ImageMagickPoint {
	return geo.ImageMagickPoint
}

type ImageMagickChannelStatistics struct {
	Min               float64 `json:"min"`               // 0,
	Max               float64 `json:"max"`               // 255,
	Mean              float64 `json:"mean"`              // 187.475,
	StandardDeviation float64 `json:"standardDeviation"` // 90.9415,
	Kurtosis          float64 `json:"kurtosis"`          // -1.22588,
	Skewness          float64 `json:"skewness"`          // -0.755169,
	Entropy           float64 `json:"entropy"`           // 0.515529
}

func (details *ImageMagickDetails) ToJSON(pretty bool) (out []byte, err error) {
	if pretty {
		return json.MarshalIndent(details, "", "  ")
	}
	return json.Marshal(details)
}

func (details *ImageMagickImageDetails) ToJSON(pretty bool) (out []byte, err error) {
	if pretty {
		return json.MarshalIndent(details, "", "  ")
	}
	return json.Marshal(details)
}

package sitemap

import (
	"encoding/xml"
)

type ImageEntry struct {
	XMLName     xml.Name `xml:"image:image"`
	Location    string   `xml:"image:loc"`
	Title       string   `xml:"image:title,omitempty"`
	Caption     string   `xml:"image:caption,omitempty"`
	GeoLocation string   `xml:"image:geo_location,omitempty"`
	License     string   `xml:"image:license,omitempty"`
}

func NewImageEntry(location, title, caption, geoLocation, license string) (ImageEntry, error) {
	return ImageEntry{
		Location:    location,
		Title:       title,
		Caption:     caption,
		GeoLocation: geoLocation,
		License:     license,
	}, nil
}

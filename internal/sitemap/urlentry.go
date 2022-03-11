package sitemap

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type ChangeFrequency string

const (
	Always  ChangeFrequency = "always"
	Hourly  ChangeFrequency = "hourly"
	Daily   ChangeFrequency = "daily"
	Weekly  ChangeFrequency = "weekly"
	Monthly ChangeFrequency = "monthly"
	Yearly  ChangeFrequency = "yearly"
	Never   ChangeFrequency = "never"
)

var allChangeFrequencies = [...]ChangeFrequency{Always, Hourly, Daily, Weekly, Monthly, Yearly, Never}

type UrlEntry struct {
	XMLName         xml.Name        `xml:"url"`
	Location        string          `xml:"loc"`
	LastModified    string          `xml:"lastmod,omitempty"`
	ChangeFrequency ChangeFrequency `xml:"changefreq,omitempty"`
	Priority        float64         `xml:"priority,omitempty"`
	Images          []ImageEntry    `xml:",omitempty"`
}

func NewUrlEntry(location, lastModified string, changeFrequency ChangeFrequency, priority float64, images []ImageEntry) (UrlEntry, error) {
	// validate Priority
	if priority < 0 || priority > 1 {
		return UrlEntry{}, errors.New("valid `priority` values range from 0.0 to 1.0")
	}

	// validate ChangeFrequency
	found := false
	for _, cf := range allChangeFrequencies {
		if changeFrequency == cf {
			found = true
		}
	}
	if !found {
		return UrlEntry{}, fmt.Errorf("valid `changeFrequency` values: '%s', '%s', '%s', '%s', '%s', '%s', '%s'", Always, Hourly, Daily, Weekly, Monthly, Yearly, Never)
	}

	// validate Images
	if len(images) > 1000 {
		return UrlEntry{}, fmt.Errorf("sitemap URL entries can only have up to 1,000 images per page")
	}

	return UrlEntry{
		Location:        location,
		LastModified:    lastModified,
		ChangeFrequency: changeFrequency,
		Priority:        priority,
		Images:          images,
	}, nil
}

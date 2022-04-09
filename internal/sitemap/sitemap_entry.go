package sitemap

import "encoding/xml"

type Entry struct {
	XMLName      xml.Name `xml:"sitemap"`
	Location     string   `xml:"loc"`
	LastModified string   `xml:"lastmod,omitempty"`
}

func NewSitemapEntry(location, lastModified string) Entry {
	return Entry{
		Location:     location,
		LastModified: lastModified,
	}
}

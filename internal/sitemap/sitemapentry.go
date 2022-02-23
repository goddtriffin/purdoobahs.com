package sitemap

import "encoding/xml"

type SitemapEntry struct {
	XMLName      xml.Name `xml:"sitemap"`
	Location     string   `xml:"loc"`
	LastModified string   `xml:"lastmod,omitempty"`
}

func NewSitemapEntry(location, lastModified string) SitemapEntry {
	return SitemapEntry{
		Location:     location,
		LastModified: lastModified,
	}
}

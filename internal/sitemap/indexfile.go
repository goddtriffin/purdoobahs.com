package sitemap

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type IndexFile struct {
	XMLName    xml.Name `xml:"sitemapindex"`
	VersionUrl string   `xml:"xmlns,attr"`
	Sitemaps   []Entry  `xml:",omitempty"`
}

func NewIndexFile(entries []Entry) *IndexFile {
	return &IndexFile{
		VersionUrl: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps:   entries,
	}
}

func (f *IndexFile) WriteToFile(filepath string) error {
	// write sitemap header line
	err := ioutil.WriteFile(filepath, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"), 0644)
	if err != nil {
		return err
	}

	// open file to append to
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	// write XML
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "    ")
	err = encoder.Encode(f)
	if err != nil {
		return err
	}

	// append newline to file
	if _, err := file.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

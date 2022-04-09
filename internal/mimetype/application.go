package mimetype

type Application MimeType

func (a Application) String() string {
	return string(a)
}

const (
	Gzip                  Application = "application/gzip"
	JavascriptApplication Application = "application/javascript"
	Json                  Application = "application/json"
	LdJson                Application = "application/ld+json"
	OggApplication        Application = "application/ogg"
	Pdf                   Application = "application/pdf"
	XTar                  Application = "application/x-tar"
	XHtmlXml              Application = "application/xhtml+xml"
	XmlApplication        Application = "application/xml"
	Zip                   Application = "application/zip"
	x7zCompressed         Application = "application/x-7z-compressed"
)

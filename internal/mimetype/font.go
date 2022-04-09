package mimetype

type Font MimeType

func (f Font) String() string {
	return string(f)
}

const (
	Otf   Font = "font/otf"
	Ttf   Font = "font/ttf"
	Woff  Font = "font/woff"
	Woff2 Font = "font/woff2"
)

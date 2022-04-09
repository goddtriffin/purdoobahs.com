package mimetype

type Image MimeType

func (i Image) String() string {
	return string(i)
}

const (
	Bmp              Image = "image/bmp"
	Gif              Image = "image/gif"
	XIcon            Image = "image/x-icon"
	VndMicrosoftIcon Image = "image/vnd.microsoft.icon"
	Jpeg             Image = "image/jpeg"
	Png              Image = "image/png"
	SvgXml           Image = "image/svg+xml"
	Tiff             Image = "image/tiff"
	Webp             Image = "image/webp"
)

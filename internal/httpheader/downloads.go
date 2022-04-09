package httpheader

type Downloads HttpHeader

func (d Downloads) String() string {
	return string(d)
}

// List of downloads HTTP headers.
const (
	ContentDisposition Downloads = "Content-Disposition"
)

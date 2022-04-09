package httpheader

type Redirects HttpHeader

func (r Redirects) String() string {
	return string(r)
}

// List of redirects HTTP headers.
const (
	Location Redirects = "Location"
)

package httpheader

type Caching HttpHeader

func (c Caching) String() string {
	return string(c)
}

// List of caching HTTP headers.
const (
	Age           Caching = "Age"
	CacheControl  Caching = "Cache-Control"
	ClearSiteData Caching = "Clear-Site-Data"
	Expires       Caching = "Expires"
	Pragma        Caching = "Pragma"
)

// List of deprecated caching HTTP headers.
const (
	DeprecatedWarning Caching = "Warning"
)

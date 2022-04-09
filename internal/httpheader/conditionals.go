package httpheader

type Conditional HttpHeader

func (c Conditional) String() string {
	return string(c)
}

// List of conditional HTTP headers.
const (
	LastModified      Conditional = "Last-Modified"
	ETag              Conditional = "ETag"
	IfMatch           Conditional = "If-Match"
	IfNoneMatch       Conditional = "If-None-Match"
	IfModifiedSince   Conditional = "If-Modified-Since"
	IfUnmodifiedSince Conditional = "If-Unmodified-Since"
	Vary              Conditional = "Vary"
)

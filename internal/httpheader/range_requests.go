package httpheader

type RangeRequests HttpHeader

func (rr RangeRequests) String() string {
	return string(rr)
}

// List of range requests HTTP headers.
const (
	AcceptRanges RangeRequests = "Accept-Ranges"
	Range        RangeRequests = "Range"
	IfRange      RangeRequests = "If-Range"
	ContentRange RangeRequests = "Content-Range"
)

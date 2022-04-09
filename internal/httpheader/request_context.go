package httpheader

type RequestContext HttpHeader

func (rc RequestContext) String() string {
	return string(rc)
}

// List of request context HTTP headers.
const (
	From           RequestContext = "From"
	Host           RequestContext = "Host"
	Referer        RequestContext = "Referer"
	ReferrerPolicy RequestContext = "Referrer-Policy"
	UserAgent      RequestContext = "User-Agent"
)

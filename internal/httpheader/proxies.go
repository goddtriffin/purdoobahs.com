package httpheader

type Proxies HttpHeader

func (p Proxies) String() string {
	return string(p)
}

// List of proxies HTTP headers.
const (
	Forwarded Proxies = "Forwarded"
	Via       Proxies = "Via"
)

// List of non-standard proxies HTTP headers.
const (
	NonstandardXForwardedFor   Proxies = "X-Forwarded-For"
	NonstandardXForwardedHost  Proxies = "X-Forwarded-Host"
	NonstandardXForwardedProto Proxies = "X-Forwarded-Proto"
)

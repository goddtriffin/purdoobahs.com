package httpheader

type Authentication HttpHeader

func (a Authentication) String() string {
	return string(a)
}

// List of authentication HTTP headers.
const (
	WwwAuthenticate    Authentication = "WWW-Authenticate"
	Authorization      Authentication = "Authorization"
	ProxyAuthenticate  Authentication = "Proxy-Authenticate"
	ProxyAuthorization Authentication = "Proxy-Authorization"
)

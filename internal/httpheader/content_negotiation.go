package httpheader

type ContentNegotiation HttpHeader

func (cn ContentNegotiation) String() string {
	return string(cn)
}

// List of content negotiation HTTP headers.
const (
	Accept         ContentNegotiation = "Accept"
	AcceptEncoding ContentNegotiation = "Accept-Encoding"
	AcceptLanguage ContentNegotiation = "Accept-Language"
)

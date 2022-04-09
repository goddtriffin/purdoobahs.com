package httpheader

type ClientHints HttpHeader

func (ch ClientHints) String() string {
	return string(ch)
}

// List of client hints HTTP headers.
const (
	AcceptCh ClientHints = "Accept-CH"
)

// List of experimental HTTP headers.
const (
	ExperimentalAcceptChLifetime ClientHints = "Accept-CH-Lifetime"
)

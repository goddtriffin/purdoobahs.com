package httpheader

type ResponseContext HttpHeader

func (rc ResponseContext) String() string {
	return string(rc)
}

// List of response context HTTP headers.
const (
	Allow  ResponseContext = "Allow"
	Server ResponseContext = "Server"
)

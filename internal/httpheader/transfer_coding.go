package httpheader

type TransferCoding HttpHeader

func (tc TransferCoding) String() string {
	return string(tc)
}

// List of transfer coding HTTP headers.
const (
	TransferEncoding TransferCoding = "Transfer-Encoding"
	Te               TransferCoding = "TE"
	Trailer          TransferCoding = "Trailer"
)

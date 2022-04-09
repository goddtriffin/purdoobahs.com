package httpheader

type NetworkClientHints HttpHeader

func (nch NetworkClientHints) String() string {
	return string(nch)
}

// List of network client hint HTTP headers.
const (
	Downlink NetworkClientHints = "Downlink"
	Ect      NetworkClientHints = "ECT"
	Rtt      NetworkClientHints = "Rtt"
)

// List of experimental network client hint HTTP headers.
const (
	SaveData NetworkClientHints = "Save-Data"
)

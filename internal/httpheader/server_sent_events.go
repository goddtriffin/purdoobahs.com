package httpheader

type ServerSentEvents HttpHeader

func (sse ServerSentEvents) String() string {
	return string(sse)
}

// List of server-sent events HTTP headers.
const (
	ReportTo ServerSentEvents = "Report-To"
)

// List of experimental server-sent events HTTP headers.
const (
	ExperimentalNel ServerSentEvents = "NEL"
)

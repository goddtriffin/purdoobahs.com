package httpheader

type HttpHeader string

func (hh HttpHeader) String() string {
	return string(hh)
}

// List of miscellaneous HTTP headers that don't have further categorization.
const (
	AltSvc               HttpHeader = "Alt-Svc"
	Date                 HttpHeader = "Date"
	Link                 HttpHeader = "Link"
	RetryAfter           HttpHeader = "Retry-After"
	ServerTiming         HttpHeader = "Server-Timing"
	ServiceWorkerAllowed HttpHeader = "Service-Worker-Allowed"
	SourceMap            HttpHeader = "SourceMap"
	Upgrade              HttpHeader = "Upgrade"
	XDnsPrefetchControl  HttpHeader = "X-DNS-Prefetch-Control"
)

// List of experimental miscellaneous HTTP headers that don't have further categorization.
const (
	ExperimentalAcceptPushPolicy HttpHeader = "Accept-Push-Policy"
	ExperimentalAcceptSignature  HttpHeader = "Accept-Signature"
	ExperimentalEarlyData        HttpHeader = "Early-Data"
	ExperimentalPushPolicy       HttpHeader = "Push-Policy"
	ExperimentalSignature        HttpHeader = "Signature"
	ExperimentalSignedHeaders    HttpHeader = "Signed-Headers"
)

// List of deprecated miscellaneous HTTP headers that don't have further categorization.
const (
	DeprecatedLargeAllocation HttpHeader = "Large-Allocation"
	DeprecatedXRobotsTag      HttpHeader = "X-Robots-Tag"
	DeprecatedXUaCompatible   HttpHeader = "X-UA-Compatible"
)

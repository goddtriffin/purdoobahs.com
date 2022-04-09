package httpheader

type HttpPublicKeyPinning HttpHeader

func (hpkp HttpPublicKeyPinning) String() string {
	return string(hpkp)
}

// List of HTTP public key pinning HTTP headers.
const (
	PublicKeyPins           HttpPublicKeyPinning = "Public-Key-Pins"
	PublicKeyPinsReportOnly HttpPublicKeyPinning = "Public-Key-Pins-Report-Only"
)

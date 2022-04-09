package httpheader

type DeviceClientHints HttpHeader

func (dch DeviceClientHints) String() string {
	return string(dch)
}

// List of deprecated device client hint HTTP headers.
const (
	DeprecatedContentDpr    DeviceClientHints = "Content-DPR"
	DeprecatedDeviceMemory  DeviceClientHints = "DeviceMemory"
	DeprecatedDpr           DeviceClientHints = "DPR"
	DeprecatedViewportWidth DeviceClientHints = "Viewport-Width"
	DeprecatedWidth         DeviceClientHints = "Width"
)

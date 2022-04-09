package httpheader

type ConnectionManagement HttpHeader

func (cm ConnectionManagement) String() string {
	return string(cm)
}

// List of connection management HTTP headers.
const (
	Connection ConnectionManagement = "Connection"
	KeepAlive  ConnectionManagement = "Keep-Alive"
)

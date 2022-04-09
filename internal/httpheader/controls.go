package httpheader

type Controls HttpHeader

func (c Controls) String() string {
	return string(c)
}

// List of controls HTTP header.
const (
	Expect Controls = "Expect"
)

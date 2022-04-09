package httpheader

type Cookies HttpHeader

func (c Cookies) String() string {
	return string(c)
}

// List of cookies HTTP headers.
const (
	Cookie    Cookies = "Cookie"
	SetCookie Cookies = "Set-Cookie"
)

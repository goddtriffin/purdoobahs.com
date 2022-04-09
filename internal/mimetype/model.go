package mimetype

type Model MimeType

func (m Model) String() string {
	return string(m)
}

const (
	Vrml Model = "model/vrml"
)

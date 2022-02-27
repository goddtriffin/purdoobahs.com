package traditions

// ITraditionService defines a Traditions Service
type ITraditionService interface {
	All() ([]*Tradition, error)
	ByName(string) (*Tradition, error)
}

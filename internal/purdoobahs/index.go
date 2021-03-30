package purdoobahs

// defines a Purdoobah Service
type IPurdoobahService interface {
	All() ([]*Purdoobah, error)
	ByName(string) (*Purdoobah, error)
	CurrentSection() (*Section, error)
	SectionByYear(int) ([]*Purdoobah, error)
}

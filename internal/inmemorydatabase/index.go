package inmemorydatabase

import "github.com/purdoobahs/purdoobahs.com/internal/purdoobahs"

type PurdoobahService struct {
	purdoobahs map[string]*purdoobahs.Purdoobah
}

func NewPurdoobahService(purdoobahs map[string]*purdoobahs.Purdoobah) *PurdoobahService {
	return &PurdoobahService{
		purdoobahs: purdoobahs,
	}
}

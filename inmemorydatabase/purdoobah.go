package inmemorydatabase

import (
	"fmt"

	"github.com/purdoobahs/purdoobahs.com/purdoobahs"
)

type PurdoobahService struct {
	purdoobahs map[string]*purdoobahs.Purdoobah
}

func NewPurdoobahService(purdoobahs map[string]*purdoobahs.Purdoobah) *PurdoobahService {
	return &PurdoobahService{
		purdoobahs: purdoobahs,
	}
}

func (ps *PurdoobahService) All() ([]*purdoobahs.Purdoobah, error) {
	allPurdoobahs := make([]*purdoobahs.Purdoobah, 0, len(ps.purdoobahs))

	for _, v := range ps.purdoobahs {
		allPurdoobahs = append(allPurdoobahs, v)
	}

	return allPurdoobahs, nil
}

func (ps *PurdoobahService) ByName(name string) (*purdoobahs.Purdoobah, error) {
	if purdoobah, ok := ps.purdoobahs[name]; ok {
		return purdoobah, nil
	}

	return &purdoobahs.Purdoobah{}, fmt.Errorf("no Purdoobah exists with that name")
}

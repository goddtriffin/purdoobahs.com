package inmemorydatabase

import (
	"fmt"
	"sort"

	"github.com/purdoobahs/purdoobahs.com/internal/traditions"
)

type TraditionService struct {
	traditions map[string]*traditions.Tradition
}

func NewTraditionService(traditions map[string]*traditions.Tradition) *TraditionService {
	return &TraditionService{
		traditions: traditions,
	}
}

// All returns every single Tradition.
func (ts *TraditionService) All() ([]*traditions.Tradition, error) {
	allTraditions := make([]*traditions.Tradition, 0, len(ts.traditions))

	for _, v := range ts.traditions {
		allTraditions = append(allTraditions, v)
	}

	sort.Sort(traditions.ByName(allTraditions))
	return allTraditions, nil
}

// ByName returns a single Tradition by their nickname.
func (ts *TraditionService) ByName(name string) (*traditions.Tradition, error) {
	if tradition, ok := ts.traditions[name]; ok {
		return tradition, nil
	}

	return &traditions.Tradition{}, fmt.Errorf("no Tradition exists with that name")
}

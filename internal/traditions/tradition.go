package traditions

import (
	"strings"
)

type Tradition struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Metadata struct {
		Image struct {
			File string `json:"file"`
			Alt  string `json:"alt"`
		} `json:"image"`
	} `json:"metadata"`
}

// ByName sorts Traditions by name
type ByName []*Tradition

func (t ByName) Len() int {
	return len(t)
}

func (t ByName) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByName) Less(i, j int) bool {
	return strings.Compare(t[i].Name, t[j].Name) < 0
}

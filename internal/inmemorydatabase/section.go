package inmemorydatabase

import (
	"sort"
	"time"

	"github.com/purdoobahs/purdoobahs.com/internal/purdoobahs"
)

func (ps *PurdoobahService) CurrentSection() (*purdoobahs.Section, error) {
	currentSection := &purdoobahs.Section{
		StudentLeaders: make([]*purdoobahs.Purdoobah, 0),
		SuperSeniors:   make([]*purdoobahs.Purdoobah, 0),
		Seniors:        make([]*purdoobahs.Purdoobah, 0),
		Juniors:        make([]*purdoobahs.Purdoobah, 0),
		Sophomores:     make([]*purdoobahs.Purdoobah, 0),
		Freshmen:       make([]*purdoobahs.Purdoobah, 0),
	}

	currentAcademicYear := ps.currentAcademicYear()

	for _, v := range ps.purdoobahs {
		if !v.MarchedDuringYear(currentAcademicYear) {
			continue
		}

		if v.IsStudentLeaderInYear(currentAcademicYear) {
			currentSection.StudentLeaders = append(currentSection.StudentLeaders, v)
		} else if v.IsYear(purdoobahs.SuperSenior) {
			currentSection.SuperSeniors = append(currentSection.SuperSeniors, v)
		} else if v.IsYear(purdoobahs.Senior) {
			currentSection.Seniors = append(currentSection.Seniors, v)
		} else if v.IsYear(purdoobahs.Junior) {
			currentSection.Juniors = append(currentSection.Juniors, v)
		} else if v.IsYear(purdoobahs.Sophomore) {
			currentSection.Sophomores = append(currentSection.Sophomores, v)
		} else if v.IsYear(purdoobahs.Freshman) {
			currentSection.Freshmen = append(currentSection.Freshmen, v)
		}
	}

	sort.Sort(purdoobahs.ByName(currentSection.StudentLeaders))
	sort.Sort(purdoobahs.ByName(currentSection.SuperSeniors))
	sort.Sort(purdoobahs.ByName(currentSection.Seniors))
	sort.Sort(purdoobahs.ByName(currentSection.Juniors))
	sort.Sort(purdoobahs.ByName(currentSection.Sophomores))
	sort.Sort(purdoobahs.ByName(currentSection.Freshmen))

	return currentSection, nil
}

func (ps *PurdoobahService) currentAcademicYear() int {
	t := time.Now()
	switch t.Month() {
	case time.January,
		time.February,
		time.March,
		time.April,
		time.May,
		time.June,
		time.July,
		time.August:
		return t.Year() - 1
	default:
		return t.Year()
	}
}

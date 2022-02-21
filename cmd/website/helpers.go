package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/purdoobahs/purdoobahs.com/internal/purdoobahs"
)

func (app *application) loadPurdoobahs() (map[string]*purdoobahs.Purdoobah, error) {
	allPurdoobahs := make(map[string]*purdoobahs.Purdoobah)

	// read in the Purdoobah JSON Schema
	filepaths, err := app.walkMatch("./assets/purdoobahs/", `*.json`)
	if err != nil {
		app.logger.Error("failed to load Purdoobah JSON filepaths")
		return allPurdoobahs, err
	}

	// loop through each file
	for _, path := range filepaths {
		// ignore _purdoobah.schema.json and _template.json
		if strings.Contains(path, "_") {
			continue
		}

		// read in the Purdoobah JSON document
		b, err := ioutil.ReadFile(path)
		if err != nil {
			app.logger.Error("failed to read in Purdoobah JSON file")
			return allPurdoobahs, err
		}

		// marshal it from JSON to struct
		var p purdoobahs.Purdoobah
		err = json.Unmarshal(b, &p)
		if err != nil {
			app.logger.Error("failed to unmarshal Purdoobah JSON")
			return allPurdoobahs, err
		}

		// generate ID (their toobah name)
		id := strings.ReplaceAll(filepath.Base(path), ".json", "")
		p.ID = id

		// generate image location
		const baseImagePath = "/static/image/purdoobah"
		if app.doesPurdoobahHaveProfilePicture(id) {
			p.Metadata.Image.File = fmt.Sprintf("%s/%s.webp", baseImagePath, id)
		} else {
			id := "_unknown"
			p.Metadata.Image.File = fmt.Sprintf("%s/%s.webp", baseImagePath, id)
		}
		p.Metadata.Image.Alt = fmt.Sprintf("%s's Profile Picture", p.Name)

		// add it to container of all purdoobahs
		allPurdoobahs[id] = &p
	}

	return allPurdoobahs, nil
}

func (app *application) walkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (app *application) doesPurdoobahHaveProfilePicture(targetID string) bool {
	// read in the Purdoobah JSON Schema
	filepaths, err := app.walkMatch("./static/image/purdoobah/", `*.webp`)
	if err != nil {
		app.logger.Error("failed to load Purdoobah image filepaths")
		os.Exit(1)
	}

	// loop through each file
	for _, path := range filepaths {
		id := strings.ReplaceAll(filepath.Base(path), ".webp", "")
		if id == targetID {
			return true
		}
	}

	app.logger.Error(fmt.Sprintf("failed to load Purdoobah image for %s", targetID))
	return false
}

func (app *application) doesSectionHaveSocialImage(targetYear int) bool {
	// read in the section social images
	filepaths, err := app.walkMatch("./static/image/section/", `*.webp`)
	if err != nil {
		app.logger.Error("failed to load section image filepaths")
		return false
	}

	// loop through each file
	for _, path := range filepaths {
		year := strings.ReplaceAll(filepath.Base(path), ".webp", "")
		yearAsInt, err := strconv.Atoi(year)
		if err != nil {
			return false
		}

		if yearAsInt == targetYear {
			return true
		}
	}

	app.logger.Error(fmt.Sprintf("failed to load section social image for year %v", targetYear))
	return false
}

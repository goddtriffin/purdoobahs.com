package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		if app.doesPurdoobahHaveProfilePicture(id) {
			p.Metadata.Image.File = fmt.Sprintf("%s.jpg", id)
		} else {
			p.Metadata.Image.File = "_unknown.jpg"
		}
		p.Metadata.Image.Alt = fmt.Sprintf("%s's Profile Picture", p.Name)

		// add it to container of all toobahs
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
	filepaths, err := app.walkMatch("./static/image/purdoobah/", `*.jpg`)
	if err != nil {
		app.logger.Error("failed to load Purdoobah image filepaths")
		os.Exit(1)
	}

	// loop through each file
	for _, path := range filepaths {
		id := strings.ReplaceAll(filepath.Base(path), ".jpg", "")
		if id == targetID {
			return true
		}
	}

	app.logger.Error(fmt.Sprintf("failed to load Purdoobah image for %s", targetID))
	return false
}

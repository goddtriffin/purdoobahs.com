package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/purdoobahs/purdoobahs.com/internal/purdoobahs"
	"github.com/purdoobahs/purdoobahs.com/internal/traditions"
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

		// generate ID (the purdoobah file name)
		id := strings.ReplaceAll(filepath.Base(path), ".json", "")
		p.ID = id

		// generate image location
		const baseImagePath = "/static/image/purdoobah"
		if app.doesPurdoobahHaveProfilePicture(id) {
			p.Metadata.Image.File = app.cacheBuster.Get(fmt.Sprintf("%s/%s.webp", baseImagePath, id))
		} else {
			id := "_unknown"
			p.Metadata.Image.File = app.cacheBuster.Get(fmt.Sprintf("%s/%s.webp", baseImagePath, id))
		}
		p.Metadata.Image.Alt = fmt.Sprintf("%s's Profile Picture", p.Name)

		// add it to container of all purdoobahs
		allPurdoobahs[id] = &p
	}

	return allPurdoobahs, nil
}

func (app *application) loadTraditions() (map[string]*traditions.Tradition, error) {
	allTraditions := make(map[string]*traditions.Tradition)

	// read in the Tradition JSON Schema
	filepaths, err := app.walkMatch("./assets/traditions/", `*.json`)
	if err != nil {
		app.logger.Error("failed to load Tradition JSON filepaths")
		return allTraditions, err
	}

	// loop through each file
	for _, path := range filepaths {
		// ignore _tradition.schema.json and _template.json
		if strings.Contains(path, "_") {
			continue
		}

		// read in the Tradition JSON document
		b, err := ioutil.ReadFile(path)
		if err != nil {
			app.logger.Error("failed to read in Tradition JSON file")
			return allTraditions, err
		}

		// marshal it from JSON to struct
		var t traditions.Tradition
		err = json.Unmarshal(b, &t)
		if err != nil {
			app.logger.Error("failed to unmarshal Tradition JSON")
			return allTraditions, err
		}

		// generate ID (the tradition file name)
		id := strings.ReplaceAll(filepath.Base(path), ".json", "")
		t.ID = id

		// generate image location
		const baseImagePath = "/static/image/tradition"
		if app.doesTraditionHavePicture(id) {
			t.Metadata.Image.File = app.cacheBuster.Get(fmt.Sprintf("%s/%s.webp", baseImagePath, id))
		} else {
			id := "_unknown"
			t.Metadata.Image.File = app.cacheBuster.Get(fmt.Sprintf("%s/%s.webp", baseImagePath, id))
		}
		t.Metadata.Image.Alt = fmt.Sprintf("%s", t.Name)

		// add it to container of all purdoobahs
		allTraditions[id] = &t
	}

	return allTraditions, nil
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

func (app *application) doesTraditionHavePicture(targetID string) bool {
	if app.cacheBuster.Get(fmt.Sprintf("/static/image/tradition/%s.webp", targetID)) == "" {
		app.logger.Error(fmt.Sprintf("tradition doesn't have an image: `%s`", targetID))
		return false
	}

	return true
}

func (app *application) doesPurdoobahHaveProfilePicture(targetID string) bool {
	if app.cacheBuster.Get(fmt.Sprintf("/static/image/purdoobah/%s.webp", targetID)) == "" {
		app.logger.Error(fmt.Sprintf("purdoobah doesn't have an image: `%s`", targetID))
		return false
	}

	return true
}

func (app *application) doesSectionHaveSocialImage(targetYear int) bool {
	if app.cacheBuster.Get(fmt.Sprintf("/static/image/section/%d.webp", targetYear)) == "" {
		app.logger.Error(fmt.Sprintf("section doesn't have an image: `%d`", targetYear))
		return false
	}

	return true
}

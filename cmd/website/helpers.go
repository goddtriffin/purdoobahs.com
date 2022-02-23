package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/purdoobahs/purdoobahs.com/internal/purdoobahs"
	"github.com/purdoobahs/purdoobahs.com/internal/sitemap"
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

func (app *application) generateIndexSitemap() error {
	homeUrl := "https://www.purdoobahs.com"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// generate index sitemap
	indexSitemap := sitemap.NewIndexFile([]sitemap.SitemapEntry{
		sitemap.NewSitemapEntry(fmt.Sprintf("%s%s", homeUrl, "/sitemap-root.xml"), lastModified),
		sitemap.NewSitemapEntry(fmt.Sprintf("%s%s", homeUrl, "/purdoobah/sitemap.xml"), lastModified),
		sitemap.NewSitemapEntry(fmt.Sprintf("%s%s", homeUrl, "/section/sitemap.xml"), lastModified),
	})
	err := indexSitemap.WriteToFile("./static/file/sitemap-index.xml")
	if err != nil {
		return err
	}

	return nil
}

func (app *application) generateRootSitemap(router *mux.Router) {
	homeUrl := "https://www.purdoobahs.com"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// gather routes for root sitemap (/sitemap-root.xml)
	rootSitemap := sitemap.NewFile([]sitemap.UrlEntry{})
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		// get path from route
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		// disallow API endpoints
		if strings.HasPrefix(pathTemplate, "/api") {
			return nil
		}

		// disallow static resources
		if strings.HasPrefix(pathTemplate, "/static") {
			return nil
		}

		// disallow profile pages (those will be done in their own sitemap-purdoobah.xml)
		if strings.HasPrefix(pathTemplate, "/purdoobah") {
			return nil
		}

		// disallow section pages (those will be done in their own sitemap-section.xml)
		if strings.HasPrefix(pathTemplate, "/section") {
			return nil
		}

		// disallow standalone files
		standalone_files := []string{"/favicon.ico", "/sitemap.xml", "/robots.txt", "/humans.txt", "/sitemap-root.xml"}
		for _, file := range standalone_files {
			if pathTemplate == file {
				return nil
			}
		}

		// add path as UrlEntry to rootSitemap
		urlEntry, err := sitemap.NewUrlEntry(fmt.Sprintf("%s%s", homeUrl, pathTemplate), lastModified, sitemap.Monthly, 0.5)
		if err != nil {
			return err
		}
		rootSitemap.AddUrl(urlEntry)

		return nil
	})
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	// generate root sitemap (GET /sitemap-root.xml)
	err = rootSitemap.WriteToFile("./static/file/sitemap-root.xml")
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

func (app *application) generateProfilesSitemap() error {
	homeUrl := "https://www.purdoobahs.com/purdoobah"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// init profiles sitemap
	profilesSitemap := sitemap.NewFile([]sitemap.UrlEntry{})

	// add new UrlEntry for every Purdoobah (by ID, not Name/BirthCertificateName)
	allPurdoobahs, err := app.purdoobahService.All()
	if err != nil {
		return err
	}
	for _, purdoobah := range allPurdoobahs {
		urlEntry, err := sitemap.NewUrlEntry(fmt.Sprintf("%s/%s", homeUrl, purdoobah.ID), lastModified, sitemap.Monthly, 0.5)
		if err != nil {
			return err
		}
		profilesSitemap.AddUrl(urlEntry)
	}

	// generate profiles sitemap (GET /purdoobah/sitemap.xml)
	err = profilesSitemap.WriteToFile("./static/file/sitemap-profiles.xml")
	if err != nil {
		return err
	}

	return nil
}

func (app *application) generateSectionsSitemap() error {
	homeUrl := "https://www.purdoobahs.com/section"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// init profiles sitemap
	profilesSitemap := sitemap.NewFile([]sitemap.UrlEntry{})

	// add new UrlEntry for every unique section year
	allSectionYears, err := app.purdoobahService.AllSectionYears()
	if err != nil {
		return err
	}
	for _, uniqueYear := range allSectionYears {
		urlEntry, err := sitemap.NewUrlEntry(fmt.Sprintf("%s/%d", homeUrl, uniqueYear), lastModified, sitemap.Monthly, 0.5)
		if err != nil {
			return err
		}
		profilesSitemap.AddUrl(urlEntry)
	}

	// generate profiles sitemap (GET /purdoobah/sitemap.xml)
	err = profilesSitemap.WriteToFile("./static/file/sitemap-sections.xml")
	if err != nil {
		return err
	}

	return nil
}

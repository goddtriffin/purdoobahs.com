package main

import (
	"fmt"
	"time"

	"github.com/purdoobahs/purdoobahs.com/internal/purdoobahs"
	"github.com/purdoobahs/purdoobahs.com/internal/sitemap"
	"github.com/purdoobahs/purdoobahs.com/internal/traditions"
)

func (app *application) generateSitemaps() error {
	// get all purdoobahs
	allPurdoobahs, err := app.purdoobahService.All()
	if err != nil {
		return err
	}

	// get all section years
	allSectionYears, err := app.purdoobahService.AllSectionYears()
	if err != nil {
		return err
	}

	// get all traditions
	allTraditions, err := app.traditionService.All()
	if err != nil {
		return err
	}

	// index
	err = generateIndexSitemap()
	if err != nil {
		return err
	}

	// root
	err = generateRootSitemap()
	if err != nil {
		return err
	}

	// profiles
	err = generateProfilesSitemap(allPurdoobahs)
	if err != nil {
		return err
	}

	// sections
	err = generateSectionsSitemap(allSectionYears)
	if err != nil {
		return err
	}

	// traditions
	err = generateTraditionsSitemap(allTraditions)
	if err != nil {
		return err
	}

	return nil
}

func generateIndexSitemap() error {
	homeUrl := "https://www.purdoobahs.com"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// generate index sitemap
	indexSitemap := sitemap.NewIndexFile([]sitemap.Entry{
		sitemap.NewSitemapEntry(fmt.Sprintf("%s%s", homeUrl, "/sitemap.xml"), lastModified),
		sitemap.NewSitemapEntry(fmt.Sprintf("%s%s", homeUrl, "/sitemap-root.xml"), lastModified),
		sitemap.NewSitemapEntry(fmt.Sprintf("%s%s", homeUrl, "/purdoobah/sitemap.xml"), lastModified),
		sitemap.NewSitemapEntry(fmt.Sprintf("%s%s", homeUrl, "/section/sitemap.xml"), lastModified),
		sitemap.NewSitemapEntry(fmt.Sprintf("%s%s", homeUrl, "/tradition/sitemap.xml"), lastModified),
	})
	err := indexSitemap.WriteToFile("./static/file/sitemap-index.xml")
	if err != nil {
		return err
	}

	return nil
}

func generateRootSitemap() error {
	homeUrl := "https://www.purdoobahs.com"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// generate root sitemap
	rootSitemap := sitemap.NewFile([]sitemap.UrlEntry{})

	// add new UrlEntry for each root route
	routes := []string{"", "alumni", "tradition", "cravers-hall-of-fame"}
	for _, route := range routes {
		urlEntry, err := sitemap.NewUrlEntry(fmt.Sprintf("%s/%s", homeUrl, route), lastModified, sitemap.Weekly, 0.5)
		if err != nil {
			return err
		}
		rootSitemap.AddUrl(urlEntry)
	}

	// generate root sitemap (GET /sitemap-root.xml)
	err := rootSitemap.WriteToFile("./static/file/sitemap-root.xml")
	if err != nil {
		return err
	}

	return nil
}

func generateProfilesSitemap(allPurdoobahs []*purdoobahs.Purdoobah) error {
	homeUrl := "https://www.purdoobahs.com/purdoobah"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// init profiles sitemap
	profilesSitemap := sitemap.NewFile([]sitemap.UrlEntry{})

	// add new UrlEntry for every Purdoobah (by ID, not Name/BirthCertificateName)
	for _, purdoobah := range allPurdoobahs {
		urlEntry, err := sitemap.NewUrlEntry(fmt.Sprintf("%s/%s", homeUrl, purdoobah.ID), lastModified, sitemap.Weekly, 0.5)
		if err != nil {
			return err
		}
		profilesSitemap.AddUrl(urlEntry)
	}

	// generate profiles sitemap (GET /purdoobah/sitemap.xml)
	err := profilesSitemap.WriteToFile("./static/file/sitemap-profiles.xml")
	if err != nil {
		return err
	}

	return nil
}

func generateSectionsSitemap(allSectionYears []int) error {
	homeUrl := "https://www.purdoobahs.com/section"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// init profiles sitemap
	profilesSitemap := sitemap.NewFile([]sitemap.UrlEntry{})

	// add new UrlEntry for every unique section year
	for _, uniqueYear := range allSectionYears {
		urlEntry, err := sitemap.NewUrlEntry(fmt.Sprintf("%s/%d", homeUrl, uniqueYear), lastModified, sitemap.Weekly, 0.5)
		if err != nil {
			return err
		}
		profilesSitemap.AddUrl(urlEntry)
	}

	// generate profiles sitemap (GET /purdoobah/sitemap.xml)
	err := profilesSitemap.WriteToFile("./static/file/sitemap-sections.xml")
	if err != nil {
		return err
	}

	return nil
}

func generateTraditionsSitemap(allTraditions []*traditions.Tradition) error {
	homeUrl := "https://www.purdoobahs.com/tradition"

	// use today's date to generate Last Modified
	lastModified := time.Now().Format(time.RFC3339)

	// init profiles sitemap
	traditionsSitemap := sitemap.NewFile([]sitemap.UrlEntry{})

	// add new UrlEntry for every Purdoobah (by ID, not Name/BirthCertificateName)
	for _, tradition := range allTraditions {
		urlEntry, err := sitemap.NewUrlEntry(fmt.Sprintf("%s/%s", homeUrl, tradition.ID), lastModified, sitemap.Weekly, 0.5)
		if err != nil {
			return err
		}
		traditionsSitemap.AddUrl(urlEntry)
	}

	// generate traditions sitemap (GET /tradition/sitemap.xml)
	err := traditionsSitemap.WriteToFile("./static/file/sitemap-traditions.xml")
	if err != nil {
		return err
	}

	return nil
}

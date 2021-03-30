package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/purdoobahs/purdoobahs.com/internal/purdoobahs"

	"github.com/MagnusFrater/fontawesome"
)

// top-level
type templateData struct {
	Metadata metadata
	Header   header
	Footer   footer
	Page     page

	PurdoobahByName *purdoobahs.Purdoobah
	Purdoobahs      []*purdoobahs.Purdoobah
	CurrentSection  *purdoobahs.Section
	Year            int
}

// layout / page / partial
type (
	metadata struct {
		LanguageCode string
		CountryCode  string
		Charset      string
		Description  string
		Project      string
		Author       string
		Twitter      twitter
		HomeURL      string
		Keywords     []string
		ThemeColor   string
		SocialImage  string
	}

	header struct {
		NavLinks    []navLink
		SocialMedia []socialMedia
	}

	footer struct {
		Copyright copyright
	}

	page struct {
		DisplayName string
		URL         string

		StyleSheets []string
		Scripts     []string
	}
)

// sub-components
type (
	twitter struct {
		Username string
	}

	copyright struct {
		Start time.Time
		End   time.Time
	}

	socialMedia struct {
		Link string
		Icon fontawesome.TemplateIcon
	}

	navLink struct {
		DisplayName string
		URL         string
	}
)

func title(s string) string {
	return strings.Title(strings.ToLower(s))
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("January 2, 2006")
}

func isoDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(time.RFC3339)
}

func keywords(keywords []string) string {
	return strings.Join(keywords, ",")
}

func subtract(x, y int) int {
	return x - y
}

func marshal(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

func prettyIntSlice(s []int) string {
	var builder string

	for i, num := range s {
		builder += strconv.Itoa(num)

		if i < len(s)-1 {
			builder += ", "
		}
	}

	return builder
}

func newTemplateCache() (map[string]*template.Template, error) {
	fa, err := fontawesome.New("./assets/icons.json")
	if err != nil {
		return nil, err
	}

	var functions = template.FuncMap{
		"capitalize":     strings.ToTitle,
		"title":          title,
		"humanDate":      humanDate,
		"isoDate":        isoDate,
		"fontawesome":    fa.SVG,
		"keywords":       keywords,
		"subtract":       subtract,
		"marshal":        marshal,
		"prettyIntSlice": prettyIntSlice,
	}

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./html/pages/*.page.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./html/layouts/*.layout.gohtml")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./html/partials/*.partial.gohtml")
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serveError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// write template to buffer to catch templating errors
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(td))
	if err != nil {
		app.serveError(w, err)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) addDefaultData(td *templateData) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.Metadata = metadata{
		LanguageCode: "en",
		CountryCode:  "US",
		Charset:      "utf-8",
		Description:  "The official website of the Purdue All-American Marching Band Toobah section.",
		Project:      "Purdoobahs",
		Author:       "Todd Everett Griffin",
		Twitter:      twitter{Username: "@goddtriffin"},
		HomeURL:      "https://www.purdoobahs.com",
		Keywords: []string{
			"purdoobahs", "purdoobah", "Purdue Toobah", "Purdue tuba",
			"Purdue", "Purdue University", "university",
			"toobah", "tuba", "sousa", "sousaphone", "helicon",
			"Orville Redenbacher",
			"Purdue All-American Marching Band", "All-American Marching Band",
			"marching band", "marching", "band",
			"YMSH", "ΨΜΣΗ",
		},
		ThemeColor: "#c28e0e",
	}

	// backup social image in case one isn't provided
	if td.Metadata.SocialImage == "" {
		td.Metadata.SocialImage = "/static/image/purdoobahs.webp"
	}

	td.Header = header{
		NavLinks: []navLink{
			{DisplayName: "Home", URL: "/"},
			{DisplayName: "F.A.Q.", URL: "/faq"},
			{DisplayName: "Cravers Hall of Fame", URL: "/cravers-hall-of-fame"},
			{DisplayName: "Alumni", URL: "/purdoobah"},
			{DisplayName: "Traditions", URL: "/traditions"},
		},
	}

	td.Footer = footer{
		Copyright: copyright{
			Start: time.Date(1889, time.May, 6, 0, 0, 0, 0, time.UTC),
			End:   time.Now(),
		},
	}

	td.Page.StyleSheets = []string{"main.css"}

	return td
}

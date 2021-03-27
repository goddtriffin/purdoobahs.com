package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/purdoobahs/purdoobahs.com/internal/inmemorydatabase"

	"github.com/purdoobahs/purdoobahs.com/internal/jsonschema"

	"github.com/purdoobahs/purdoobahs.com/internal/purdoobahs"

	"github.com/MagnusFrater/helmet"
)

type environment int

const (
	development environment = iota
	production
)

type application struct {
	env           environment
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	helmet        *helmet.Helmet

	purdoobahService purdoobahs.IPurdoobahService
}

func main() {
	// parse command line flags
	addr := flag.String("addr", "", "HTTP network address")
	env := flag.String("env", "", "dictates application environment")
	flag.Parse()

	// initialize the application
	app := &application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		helmet:   createHelmet(),
	}

	// validate all Purdoobah JSON schema files
	invalidFiles, err := jsonschema.ValidateJsonSchema(app.infoLog, app.errorLog)
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	if invalidFiles {
		app.errorLog.Fatalln("Invalid Purdoobah JSON detected - exiting.")
	}

	// load Purdoobah files into Purdoobah service
	allPurdoobahs, err := loadPurdoobahs()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	app.purdoobahService = inmemorydatabase.NewPurdoobahService(allPurdoobahs)

	// switch port to serve on based on environment deployed in
	switch strings.ToLower(*env) {
	case "dev", "develop", "development":
		app.env = development

		if *addr == "" {
			*addr = ":443"
		}
	case "prod", "production":
		app.env = production
		if *addr == "" {
			*addr = ":80"
		}
	default:
		app.errorLog.Println("-env flag needs to be one of: 'dev', 'develop', 'development', 'prod', or 'production'")
		flag.Usage()
		os.Exit(1)
	}

	// create HTML template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	app.templateCache = templateCache

	// create the server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// if this is the local environment, load dummy TLS files
	if app.env == development {
		srv.TLSConfig = &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		}
	}

	// start the server
	app.infoLog.Printf("Starting server on %s\n", *addr)
	switch app.env {
	case development:
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	case production:
		err = srv.ListenAndServe()
	}

	// print error on exit
	app.errorLog.Fatalln(err)
}

func loadPurdoobahs() (map[string]*purdoobahs.Purdoobah, error) {
	allPurdoobahs := make(map[string]*purdoobahs.Purdoobah)

	// read in the Purdoobah JSON Schema
	filepaths, err := walkMatch("./assets/purdoobahs/", `*.json`)
	if err != nil {
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
			return allPurdoobahs, err
		}

		// marshal it from JSON to struct
		var p purdoobahs.Purdoobah
		err = json.Unmarshal(b, &p)
		if err != nil {
			return allPurdoobahs, err
		}

		// generate ID (their toobah name)
		id := strings.ReplaceAll(filepath.Base(path), ".json", "")
		p.ID = id

		// generate image location
		if doesPurdoobahHaveProfilePicture(id) {
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

func walkMatch(root, pattern string) ([]string, error) {
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

func doesPurdoobahHaveProfilePicture(id string) bool {
	switch id {
	case "brave-little-toaster",
		"bucket",
		"cowboy",
		"crabcakes",
		"domino",
		"guido",
		"juice",
		"manbearpig",
		"mr-moeschberger",
		"ocarina",
		"pfreys",
		"poppin-fresh",
		"professor-x",
		"remington",
		"shirt",
		"stark",
		"trumoo",
		"velveeta":
		return false
	default:
		return true
	}
}

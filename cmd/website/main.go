package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/purdoobahs/purdoobahs.com/inmemorydatabase"

	"github.com/purdoobahs/purdoobahs.com/jsonschema"

	"github.com/purdoobahs/purdoobahs.com/purdoobahs"

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
	templateCache, err := newTemplateCache("./ui/")
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

	return allPurdoobahs, nil
}

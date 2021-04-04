package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/purdoobahs/purdoobahs.com/internal/logger"

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
	logger        logger.ILogger
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
		logger: logger.NewLogger(),
		helmet: createHelmet(),
	}

	// validate all Purdoobah JSON schema files
	invalidFiles, err := jsonschema.ValidateJsonSchema(app.logger)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	if invalidFiles {
		app.logger.Error("Invalid Purdoobah JSON detected - exiting.")
		os.Exit(1)
	}

	// load Purdoobah files into Purdoobah service
	allPurdoobahs, err := app.loadPurdoobahs()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
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
		app.logger.Error("-env flag needs to be one of: 'dev', 'develop', 'development', 'prod', or 'production'")
		flag.Usage()
		os.Exit(1)
	}

	// create HTML template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		app.logger.Error(err.Error())
	}
	app.templateCache = templateCache

	// create the server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.logger.(*logger.Logger).ErrorLog,
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
	app.logger.Info(fmt.Sprintf("Starting server on %s\n", *addr))
	switch app.env {
	case development:
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	case production:
		err = srv.ListenAndServe()
	}

	// print error on exit
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

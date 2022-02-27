package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/purdoobahs/purdoobahs.com/internal/logger"
	"github.com/purdoobahs/purdoobahs.com/internal/traditions"

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
	traditionService traditions.ITraditionService

	httpClient *http.Client
}

func main() {
	// parse command line flags
	addr := flag.String("addr", "", "HTTP network address")
	env := flag.String("env", "", "dictates application environment")
	flag.Parse()

	// set default address if it isn't set
	if *addr == "" {
		*addr = ":8080"
	}

	// initialize the application
	app := &application{
		logger: logger.NewLogger(),
		helmet: createHelmet(),
	}

	// generate index/root sitemaps
	err := app.generateIndexSitemap()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	err = app.generateRootSitemap()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	// validate all Purdoobah JSON schema files
	invalidFiles, err := jsonschema.ValidateJsonSchema(app.logger)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	if invalidFiles {
		app.logger.Error("Invalid JSON Schema detected - exiting.")
		os.Exit(1)
	}

	// load Purdoobah files into Purdoobah service
	allPurdoobahs, err := app.loadPurdoobahs()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	app.purdoobahService = inmemorydatabase.NewPurdoobahService(allPurdoobahs)

	// generate profiles/sections sitemaps
	err = app.generateProfilesSitemap()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	err = app.generateSectionsSitemap()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	// load Tradition files into Tradition service
	allTraditions, err := app.loadTraditions()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	app.traditionService = inmemorydatabase.NewTraditionService(allTraditions)

	// generate traditions sitemaps
	err = app.generateTraditionsSitemap()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	// set environment
	switch strings.ToLower(*env) {
	case "dev", "develop", "development":
		app.env = development
	case "prod", "production":
		app.env = production
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

	// create http Client for Analytics API
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.MaxIdleConns = 100
	tr.MaxConnsPerHost = 100
	tr.MaxIdleConnsPerHost = 100
	app.httpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
	}

	// create the server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.logger.(*logger.Logger).ErrorLog,
		Handler:  app.routes(),

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start the server
	app.logger.Info(fmt.Sprintf("Starting server on %s", *addr))
	err = srv.ListenAndServe()

	// print error on exit
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

package main

import (
	"crypto/tls"
	"flag"
	"github.com/MagnusFrater/helmet"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
}

func main() {
	addr := flag.String("addr", "", "HTTP network address")
	env := flag.String("env", "", "dictates application environment")
	flag.Parse()

	app := &application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		helmet:   createHelmet(),
	}

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
		log.Println("-env flag needs to be one of: 'dev', 'develop', 'development', 'prod', or 'production'")
		flag.Usage()
		os.Exit(1)
	}

	templateCache, err := newTemplateCache("./ui/")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.templateCache = templateCache

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if app.env == development {
		srv.TLSConfig = &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		}
	}

	app.infoLog.Printf("Starting server on %s\n", *addr)

	switch app.env {
	case development:
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	case production:
		err = srv.ListenAndServe()
	}

	app.errorLog.Fatal(err)
}

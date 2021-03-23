package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.helmet.Secure)

	mux := http.NewServeMux()
	mux.HandleFunc("/faq", app.faq)
	mux.HandleFunc("/cravers-hall-of-fame", app.craversHallOfFame)
	mux.HandleFunc("/alumni", app.alumni)

	mux.HandleFunc("/favicon.ico", app.favicon)
	mux.HandleFunc("/robots.txt", app.robotsTxt)
	mux.HandleFunc("/humans.txt", app.humansTxt)

	mux.HandleFunc("/health-check", app.healthCheck)

	mux.HandleFunc("/", app.home)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

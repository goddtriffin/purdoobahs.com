package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.helmet.Secure)

	router := mux.NewRouter()
	apiSubrouter := router.PathPrefix("/api").Subrouter()
	apiPurdoobahSubrouter := apiSubrouter.PathPrefix("/purdoobah").Subrouter()

	router.HandleFunc("/faq", app.faq).Methods("GET")
	router.HandleFunc("/cravers-hall-of-fame", app.craversHallOfFame).Methods("GET")
	router.HandleFunc("/alumni", app.alumni).Methods("GET")
	router.HandleFunc("/traditions", app.traditions).Methods("GET")
	router.HandleFunc("/purdoobah/{name}", app.purdoobahProfile).Methods("GET")
	router.HandleFunc("/favicon.ico", app.favicon).Methods("GET")
	router.HandleFunc("/robots.txt", app.robotsTxt).Methods("GET")
	router.HandleFunc("/humans.txt", app.humansTxt).Methods("GET")
	router.HandleFunc("/health-check", app.healthCheck).Methods("GET")
	router.Handle("/static/", http.StripPrefix(
		"/static",
		http.FileServer(http.Dir("./ui/static")),
	),
	).Methods("GET")
	router.HandleFunc("/", app.home).Methods("GET")

	apiPurdoobahSubrouter.HandleFunc("/all", app.allPurdoobahs).Methods("GET")
	apiPurdoobahSubrouter.HandleFunc("/{name}", app.purdoobahByName).Methods("GET")

	return standardMiddleware.Then(router)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Page: page{
			DisplayName: "Home",
			URL:         "/",
		},
	})
}

func (app *application) faq(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "faq.page.tmpl", &templateData{
		Page: page{
			DisplayName: "F.A.Q.",
			URL:         "/faq",
		},
	})
}

func (app *application) craversHallOfFame(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "cravers-hall-of-fame.page.tmpl", &templateData{
		Page: page{
			DisplayName: "Cravers Hall of Fame",
			URL:         "/cravers-hall-of-fame",
		},
	})
}

func (app *application) alumni(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "alumni.page.tmpl", &templateData{
		Page: page{
			DisplayName: "Alumni",
			URL:         "/alumni",
		},
	})
}

func (app *application) traditions(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "traditions.page.tmpl", &templateData{
		Page: page{
			DisplayName: "Traditions",
			URL:         "/traditions",
		},
	})
}

func (app *application) purdoobahProfile(w http.ResponseWriter, r *http.Request) {
	// get name
	vars := mux.Vars(r)
	name := vars["name"]

	// get purdoobah
	purdoobahByName, err := app.purdoobahService.ByName(name)
	if err != nil {
		app.notFound(w)
		return
	}

	app.render(w, r, "purdoobah-profile.page.tmpl", &templateData{
		Page: page{
			DisplayName: purdoobahByName.Name,
			URL:         fmt.Sprintf("/purdoobah/%s", name),
		},
		PurdoobahByName: purdoobahByName,
	})
}

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/image/favicon.ico", http.StatusMovedPermanently)
}

func (app *application) robotsTxt(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/file/robots.txt", http.StatusMovedPermanently)
}

func (app *application) humansTxt(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/file/humans.txt", http.StatusMovedPermanently)
}

func (app *application) allPurdoobahs(w http.ResponseWriter, r *http.Request) {
	// get all purdoobahs
	allPurdoobahs, err := app.purdoobahService.All()
	if err != nil {
		app.serveError(w, err)
		return
	}

	// convert to JSON bytes
	b, err := json.Marshal(allPurdoobahs)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// send it out
	_, err = w.Write(b)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) purdoobahByName(w http.ResponseWriter, r *http.Request) {
	// get name
	vars := mux.Vars(r)
	name := vars["name"]

	// get purdoobah
	purdoobahByName, err := app.purdoobahService.ByName(name)
	if err != nil {
		app.notFound(w)
		return
	}

	// convert to JSON bytes
	b, err := json.Marshal(purdoobahByName)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// send it out
	_, err = w.Write(b)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

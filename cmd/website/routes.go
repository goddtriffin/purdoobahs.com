package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.helmet.Secure)

	mux := http.NewServeMux()
	mux.HandleFunc("/faq", app.faq)
	mux.HandleFunc("/cravers-hall-of-fame", app.craversHallOfFame)
	mux.HandleFunc("/alumni", app.alumni)
	mux.HandleFunc("/traditions", app.traditions)

	mux.HandleFunc("/favicon.ico", app.favicon)
	mux.HandleFunc("/robots.txt", app.robotsTxt)
	mux.HandleFunc("/humans.txt", app.humansTxt)

	mux.HandleFunc("/health-check", app.healthCheck)

	mux.HandleFunc("/", app.home)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
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
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, r, "faq.page.tmpl", &templateData{
		Page: page{
			DisplayName: "F.A.Q.",
			URL:         "/faq",
		},
	})
}

func (app *application) craversHallOfFame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, r, "cravers-hall-of-fame.page.tmpl", &templateData{
		Page: page{
			DisplayName: "Cravers Hall of Fame",
			URL:         "/cravers-hall-of-fame",
		},
	})
}

func (app *application) alumni(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, r, "alumni.page.tmpl", &templateData{
		Page: page{
			DisplayName: "Alumni",
			URL:         "/alumni",
		},
	})
}

func (app *application) traditions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, r, "traditions.page.tmpl", &templateData{
		Page: page{
			DisplayName: "Traditions",
			URL:         "/traditions",
		},
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

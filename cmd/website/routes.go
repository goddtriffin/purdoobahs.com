package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/purdoobahs/purdoobahs.com/internal/logger"
	"github.com/purdoobahs/purdoobahs.com/internal/plausibleanalytics"

	"github.com/gorilla/mux"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.helmet.Secure)

	// routers
	router := mux.NewRouter()
	apiSubrouter := router.PathPrefix("/api/v1").Subrouter()
	apiPurdoobahSubrouter := apiSubrouter.PathPrefix("/purdoobah").Subrouter()
	apiSectionSubrouter := apiSubrouter.PathPrefix("/section").Subrouter()

	// pages
	router.HandleFunc("/cravers-hall-of-fame", app.pageCraversHallOfFame).Methods("GET")
	router.HandleFunc("/traditions", app.pageTraditions).Methods("GET")
	router.HandleFunc("/purdoobah/{name}", app.pagePurdoobahProfile).Methods("GET")
	router.HandleFunc("/alumni", app.pageAlumni).Methods("GET")
	router.HandleFunc("/section/{year}", app.pageSectionByYear).Methods("GET")

	// files
	router.HandleFunc("/favicon.ico", app.fileFavicon).Methods("GET")
	router.HandleFunc("/robots.txt", app.fileRobotsTxt).Methods("GET")
	router.HandleFunc("/humans.txt", app.fileHumansTxt).Methods("GET")

	// static files
	router.PathPrefix("/static/").
		Handler(http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static")),
		))

	// has to occur last because it is the most generic route "/"
	router.HandleFunc("/", app.pageHome).Methods("GET")

	// generic API
	apiSubrouter.HandleFunc("/health-check", app.apiHealthCheck).Methods("GET")

	// analytics API
	apiSubrouter.HandleFunc("/scitylana", app.apiAnalytics).Methods("POST")

	// Purdoobah API
	apiPurdoobahSubrouter.HandleFunc("/all", app.apiAllPurdoobahs).Methods("GET")
	apiPurdoobahSubrouter.HandleFunc("/{name}", app.apiPurdoobahByName).Methods("GET")

	// section API
	apiSectionSubrouter.HandleFunc("/current", app.apiCurrentSection).Methods("GET")
	apiSectionSubrouter.HandleFunc("/{year}", app.apiSectionByYear).Methods("GET")

	return standardMiddleware.Then(router)
}

func (app *application) pageHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// get current section
	currentSection, err := app.purdoobahService.CurrentSection()
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "home.gohtml", &templateData{
		Page: page{
			DisplayName: "Home",
			URL:         "/",
		},
		CurrentSection: currentSection,
	})
}

func (app *application) pageCraversHallOfFame(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "cravers-hall-of-fame.gohtml", &templateData{
		Page: page{
			DisplayName: "Cravers Hall of Fame",
			URL:         "/cravers-hall-of-fame",
		},
		Metadata: metadata{
			SocialImage: "/static/image/socials/cravers-hall-of-fame.webp",
			Description: "Inductees of the 2019 White Castle Cravers Hall of Fame!",
		},
	})
}

func (app *application) pageTraditions(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "traditions.gohtml", &templateData{
		Page: page{
			DisplayName: "Traditions",
			URL:         "/traditions",
		},
	})
}

func (app *application) pagePurdoobahProfile(w http.ResponseWriter, r *http.Request) {
	// get name
	vars := mux.Vars(r)
	name := vars["name"]

	// get purdoobah
	purdoobahByName, err := app.purdoobahService.ByName(name)
	if err != nil {
		app.notFound(w)
		return
	}

	app.render(w, r, "purdoobah-profile.gohtml", &templateData{
		Page: page{
			DisplayName: fmt.Sprintf("%s %s", purdoobahByName.Name, purdoobahByName.Emoji),
			URL:         fmt.Sprintf("/purdoobah/%s", name),
		},
		PurdoobahByName: purdoobahByName,
		Metadata: metadata{
			SocialImage: purdoobahByName.Metadata.Image.File,
			Description: fmt.Sprintf(
				"Meet %s! %s Member of the %s Purdoobah section(s).",
				purdoobahByName.Name,
				purdoobahByName.Emoji,
				prettyIntSlice(purdoobahByName.Marching.YearsMarched),
			),
		},
	})
}

func (app *application) pageAlumni(w http.ResponseWriter, r *http.Request) {
	// get all purdoobahs
	allPurdoobahs, err := app.purdoobahService.All()
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "alumni.gohtml", &templateData{
		Page: page{
			DisplayName: "Alumni",
			URL:         "/alumni",
			Scripts:     []string{"alumni.js"},
		},
		Purdoobahs: allPurdoobahs,
		Metadata: metadata{
			SocialImage: "/static/image/section/2019.webp",
			Description: "OOOOOOOOOOOOOOOOOOLLLLDDDDDD",
		},
	})
}

func (app *application) pageSectionByYear(w http.ResponseWriter, r *http.Request) {
	// get year
	vars := mux.Vars(r)
	yearAsString := vars["year"]

	// convert from string to int
	yearAsInt, err := strconv.Atoi(yearAsString)
	if err != nil {
		app.notFound(w)
		return
	}

	// get section by year
	sectionByYear, err := app.purdoobahService.SectionByYear(yearAsInt)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// get social image
	socialImage := fmt.Sprintf("/static/image/section/%v.webp", yearAsInt)
	if !app.doesSectionHaveSocialImage(yearAsInt) {
		socialImage = ""
	}

	app.render(w, r, "section-by-year.gohtml", &templateData{
		Page: page{
			DisplayName: fmt.Sprintf("%d Section", yearAsInt),
			URL:         fmt.Sprintf("/section/%d", yearAsInt),
		},
		Purdoobahs: sectionByYear,
		Year:       yearAsInt,
		Metadata: metadata{
			SocialImage: socialImage,
			Description: "OOOOOOOOOOOOOOOOOOLLLLDDDDDD",
		},
	})
}

func (app *application) fileFavicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/image/favicon.ico", http.StatusMovedPermanently)
}

func (app *application) fileRobotsTxt(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/file/robots.txt", http.StatusMovedPermanently)
}

func (app *application) fileHumansTxt(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/file/humans.txt", http.StatusMovedPermanently)
}

func (app *application) apiHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) apiAnalytics(w http.ResponseWriter, r *http.Request) {
	// body
	screenWidth, err := strconv.Atoi(r.FormValue("screen_width"))
	if err != nil {
		app.serveError(w, err)
		return
	}
	body := plausibleanalytics.NewPlausibleAnalyticsBody(r.FormValue("url"), r.FormValue("referrer"), screenWidth)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// create request
	req, err := http.NewRequestWithContext(r.Context(), "POST", "https://plausible.io/api/event", bytes.NewBuffer(bodyBytes))
	if err != nil {
		app.serveError(w, err)
		return
	}

	// headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", r.FormValue("user_agent"))
	req.Header.Set("X-Forwarded-For", r.RemoteAddr)

	// print headers
	if reqHeadersBytes, err := json.Marshal(req.Header); err == nil {
		app.logger.Info(fmt.Sprintf("Plausible Analytics headers: %v", string(reqHeadersBytes)))
	}

	// print body
	app.logger.Info(fmt.Sprintf("Plausible Analytics body: %v", string(bodyBytes)))

	if app.env == production {
		// POST analytics event
		resp, err := app.httpClient.Do(req)
		if err != nil {
			app.serveError(w, err)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		// print response
		app.logger.Info(fmt.Sprintf("Plausible Analytics status: %v %v", resp.Status, resp.Header))
		app.logger.Info(fmt.Sprintf("Plausible Analytics body: %v", string(body)))
	} else {
		app.logger.Info("Not sending Plausible analytics request due to being in development environment.")
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) apiAllPurdoobahs(w http.ResponseWriter, r *http.Request) {
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

func (app *application) apiPurdoobahByName(w http.ResponseWriter, r *http.Request) {
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

func (app *application) apiCurrentSection(w http.ResponseWriter, r *http.Request) {
	// get current section
	currentSection, err := app.purdoobahService.CurrentSection()
	if err != nil {
		app.notFound(w)
		return
	}

	// convert to JSON bytes
	b, err := json.Marshal(currentSection)
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

func (app *application) apiSectionByYear(w http.ResponseWriter, r *http.Request) {
	// get year
	vars := mux.Vars(r)
	yearAsString := vars["year"]

	// convert from string to int
	yearAsInt, err := strconv.Atoi(yearAsString)
	if err != nil {
		app.notFound(w)
		return
	}

	// get section by year
	sectionByYear, err := app.purdoobahService.SectionByYear(yearAsInt)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// convert to JSON bytes
	b, err := json.Marshal(sectionByYear)
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
	err = app.logger.(*logger.Logger).ErrorLog.Output(2, trace)
	if err != nil {
		app.serveError(w, err)
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

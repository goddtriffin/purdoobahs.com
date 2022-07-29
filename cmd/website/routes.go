package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/purdoobahs/purdoobahs.com/internal/httpheader"
	"github.com/purdoobahs/purdoobahs.com/internal/logger"
	"github.com/purdoobahs/purdoobahs.com/internal/mimetype"
	"github.com/purdoobahs/purdoobahs.com/internal/plausibleanalytics"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(
		app.recoverPanic,
		handlers.ProxyHeaders,
		app.logRequest,
		app.helmet.Secure,
		app.cacheControl.ForeverCache,
	)

	// routers
	router := mux.NewRouter()
	staticFilesSubrouter := router.PathPrefix("/static").Subrouter()
	apiSubrouter := router.PathPrefix("/api").Subrouter()
	apiV1Subrouter := apiSubrouter.PathPrefix("/v1").Subrouter()

	// files
	router.HandleFunc("/favicon.ico", app.fileFavicon).Methods("GET")
	router.HandleFunc("/sitemap.xml", app.fileIndexSitemapXml).Methods("GET")
	router.HandleFunc("/sitemap-root.xml", app.fileRootSitemapXml).Methods("GET")
	router.HandleFunc("/purdoobah/sitemap.xml", app.fileProfilesSitemapXml).Methods("GET")
	router.HandleFunc("/section/sitemap.xml", app.fileSectionsSitemapXml).Methods("GET")
	router.HandleFunc("/tradition/sitemap.xml", app.fileTraditionsSitemapXml).Methods("GET")
	router.HandleFunc("/robots.txt", app.fileRobotsTxt).Methods("GET")
	router.HandleFunc("/humans.txt", app.fileHumansTxt).Methods("GET")

	// pages
	router.HandleFunc("/cravers-hall-of-fame", app.pageCraversHallOfFame).Methods("GET")
	router.HandleFunc("/tradition", app.pageTradition).Methods("GET")
	router.HandleFunc("/tradition/{name}", app.pageTraditionProfile).Methods("GET")
	router.HandleFunc("/alumni", app.pageAlumni).Methods("GET")
	router.HandleFunc("/section/{year}", app.pageSectionByYear).Methods("GET")
	router.HandleFunc("/purdoobah/{name}", app.pagePurdoobahProfile).Methods("GET")

	// static files
	staticFilesSubrouter.PathPrefix("/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	)

	// catch all
	// has to occur last because it is the most generic route "/"
	router.PathPrefix("/").HandlerFunc(app.pageHome).Methods("GET")

	// generic API
	apiV1Subrouter.HandleFunc("/health", app.apiHealthCheck).Methods("GET")

	// analytics API
	apiV1Subrouter.HandleFunc("/scitylana", app.apiAnalytics).Methods("POST")

	// Purdoobah API
	apiV1PurdoobahSubrouter := apiV1Subrouter.PathPrefix("/purdoobah").Subrouter()
	apiV1PurdoobahSubrouter.HandleFunc("/all", app.apiAllPurdoobahs).Methods("GET")
	apiV1PurdoobahSubrouter.HandleFunc("/{name}", app.apiPurdoobahByName).Methods("GET")

	// section API
	apiV1SectionSubrouter := apiV1Subrouter.PathPrefix("/section").Subrouter()
	apiV1SectionSubrouter.HandleFunc("/current", app.apiCurrentSection).Methods("GET")
	apiV1SectionSubrouter.HandleFunc("/{year}", app.apiSectionByYear).Methods("GET")

	// tradition API
	apiV1TraditionSubrouter := apiV1Subrouter.PathPrefix("/tradition").Subrouter()
	apiV1TraditionSubrouter.HandleFunc("/all", app.apiAllTraditions).Methods("GET")

	// api catch all
	// has to occur last because it is the most generic route "/api/" and "/api/v1/"
	apiSubrouter.PathPrefix("/").HandlerFunc(app.apiNotFound).Methods("GET")
	apiSubrouter.PathPrefix("").HandlerFunc(app.apiNotFound).Methods("GET")

	return standardMiddleware.Then(router)
}

func (app *application) pageHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.pageNotFound(w, r)
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
			SocialImage: app.cacheBuster.Get("/static/image/socials/cravers-hall-of-fame.webp"),
			Description: "Inductees of the 2019 White Castle Cravers Hall of Fame!",
		},
	})
}

func (app *application) pageTradition(w http.ResponseWriter, r *http.Request) {
	// get all traditions
	allTraditions, err := app.traditionService.All()
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "tradition.gohtml", &templateData{
		Page: page{
			DisplayName: "Traditions",
			URL:         "/tradition",
		},
		Traditions: allTraditions,
		Metadata: metadata{
			SocialImage: app.cacheBuster.Get("/static/image/socials/traditions.webp"),
			Description: "It's surprising what you get when you put a bunch of toobahs together in the same room.",
		},
	})
}

func (app *application) pageTraditionProfile(w http.ResponseWriter, r *http.Request) {
	// get name
	vars := mux.Vars(r)
	name := vars["name"]

	// get tradition
	traditionByName, err := app.traditionService.ByName(name)
	if err != nil {
		app.pageNotFound(w, r)
		return
	}

	app.render(w, r, "tradition-profile.gohtml", &templateData{
		Page: page{
			DisplayName: traditionByName.Name,
			URL:         fmt.Sprintf("/tradition/%s", name),
		},
		TraditionByName: traditionByName,
		Metadata: metadata{
			SocialImage: traditionByName.Metadata.Image.File,
			Description: traditionByName.Description,
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
		app.pageNotFound(w, r)
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

	// get all years marched
	allYearsMarched, err := app.purdoobahService.AllSectionYears()
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "alumni.gohtml", &templateData{
		Page: page{
			DisplayName: "Alumni",
			URL:         "/alumni",
			Scripts:     []string{app.cacheBuster.Get("/static/script/alumni.js")},
		},
		Purdoobahs:      allPurdoobahs,
		AllYearsMarched: allYearsMarched,
		Metadata: metadata{
			SocialImage: app.cacheBuster.Get("/static/image/section/2019.webp"),
			Description: "OOOOOOOOOOOOOOOOOOLLLLDDDDDD",
		},
	})
}

func (app *application) pageSectionByYear(w http.ResponseWriter, r *http.Request) {
	// get year
	vars := mux.Vars(r)
	yearAsString := vars["year"]

	// convert from string to int
	// (section -1 is for purdoobahs we don't know their marching history)
	yearAsInt := -1
	var err error
	if yearAsString != "unknown" {
		yearAsInt, err = strconv.Atoi(yearAsString)
		if err != nil {
			app.pageNotFound(w, r)
			return
		}
	}

	// get section by year
	sectionByYear, err := app.purdoobahService.SectionByYear(yearAsInt)
	if err != nil {
		app.serveError(w, err)
		return
	}

	if len(sectionByYear) == 0 {
		app.pageNotFound(w, r)
		return
	}

	// get social image
	var socialImage string
	if yearAsInt == -1 {
		socialImage = app.cacheBuster.Get("/static/image/section/unknown.webp")
	} else {
		socialImage = app.cacheBuster.Get(fmt.Sprintf("/static/image/section/%v.webp", yearAsInt))
	}
	if !app.doesSectionHaveSocialImage(yearAsInt) {
		socialImage = ""
	}

	var displayName string
	if yearAsInt == -1 {
		displayName = "Unknown Section"
	} else {
		displayName = fmt.Sprintf("%d Section", yearAsInt)
	}

	var url string
	if yearAsInt == -1 {
		url = "/section/unknown"
	} else {
		url = fmt.Sprintf("/section/%d", yearAsInt)
	}

	// get all years marched
	allYearsMarched, err := app.purdoobahService.AllSectionYears()
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "section-by-year.gohtml", &templateData{
		Page: page{
			DisplayName: displayName,
			URL:         url,
		},
		Purdoobahs:      sectionByYear,
		Year:            yearAsInt,
		AllYearsMarched: allYearsMarched,
		Metadata: metadata{
			SocialImage: socialImage,
			Description: "OOOOOOOOOOOOOOOOOOLLLLDDDDDD",
		},
	})
}

func (app *application) pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	app.render(w, r, "404.gohtml", &templateData{
		Page: page{
			DisplayName: "404",
			URL:         r.URL.Path,
		},
		Metadata: metadata{
			SocialImage: app.cacheBuster.Get("/static/image/socials/404.webp"),
		},
	})
}

func (app *application) fileFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(httpheader.ContentType.String(), mimetype.XIcon.String())
	http.ServeFile(w, r, fmt.Sprintf(".%s", app.cacheBuster.Get("/static/image/favicon/favicon.ico")))
}

func (app *application) fileIndexSitemapXml(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(httpheader.ContentType.String(), mimetype.XmlApplication.String())
	http.ServeFile(w, r, fmt.Sprintf(".%s", app.cacheBuster.Get("/static/file/sitemap-index.xml")))
}

func (app *application) fileRootSitemapXml(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(httpheader.ContentType.String(), mimetype.XmlApplication.String())
	http.ServeFile(w, r, fmt.Sprintf(".%s", app.cacheBuster.Get("/static/file/sitemap-root.xml")))
}

func (app *application) fileProfilesSitemapXml(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(httpheader.ContentType.String(), mimetype.XmlApplication.String())
	http.ServeFile(w, r, fmt.Sprintf(".%s", app.cacheBuster.Get("/static/file/sitemap-profiles.xml")))
}

func (app *application) fileSectionsSitemapXml(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(httpheader.ContentType.String(), mimetype.XmlApplication.String())
	http.ServeFile(w, r, fmt.Sprintf(".%s", app.cacheBuster.Get("/static/file/sitemap-sections.xml")))
}

func (app *application) fileTraditionsSitemapXml(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(httpheader.ContentType.String(), mimetype.XmlApplication.String())
	http.ServeFile(w, r, fmt.Sprintf(".%s", app.cacheBuster.Get("/static/file/sitemap-traditions.xml")))
}

func (app *application) fileRobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(httpheader.ContentType.String(), mimetype.Plain.String())
	http.ServeFile(w, r, fmt.Sprintf(".%s", app.cacheBuster.Get("/static/file/robots.txt")))
}

func (app *application) fileHumansTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(httpheader.ContentType.String(), mimetype.Plain.String())
	http.ServeFile(w, r, fmt.Sprintf(".%s", app.cacheBuster.Get("/static/file/humans.txt")))
}

func (app *application) apiHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(
		httpheader.ContentType.String(),
		fmt.Sprintf("%s; charset=utf-8", mimetype.Json.String()),
	)
	_, err := w.Write([]byte("{ \"status\": \"success\"}"))
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) apiAnalytics(w http.ResponseWriter, r *http.Request) {
	// set domain depending on if we're in the dev or prod environment
	domain := ""
	if app.env == production {
		domain = "purdoobahs.com"
	} else {
		domain = "test.toddgriffin.me"
	}

	// body
	screenWidth, err := strconv.Atoi(r.FormValue("screen_width"))
	if err != nil {
		app.serveError(w, err)
		return
	}
	body := plausibleanalytics.NewPlausibleAnalyticsBody(
		domain,
		"pageview",
		r.FormValue("url"),
		r.FormValue("referrer"),
		screenWidth,
	)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// create request
	req, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		"https://plausible.io/api/event",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// headers
	req.Header.Set(httpheader.ContentType.String(), mimetype.Json.String())
	req.Header.Set(httpheader.UserAgent.String(), r.FormValue("user_agent"))
	req.Header.Set(httpheader.NonstandardXForwardedFor.String(), r.RemoteAddr)

	// print headers
	if reqHeadersBytes, err := json.Marshal(req.Header); err == nil {
		app.logger.Info(fmt.Sprintf("Plausible Analytics headers: %v", string(reqHeadersBytes)))
	}

	// print body
	app.logger.Info(fmt.Sprintf("Plausible Analytics body: %v", string(bodyBytes)))

	// POST analytics event
	resp, err := app.httpClient.Do(req)
	if err != nil {
		app.serveError(w, err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			app.logger.Error(err.Error())
		}
	}()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	// print response
	app.logger.Info(fmt.Sprintf("Plausible Analytics status: %v %v", resp.Status, resp.Header))
	app.logger.Info(fmt.Sprintf("Plausible Analytics body: %v", string(resp_body)))

	w.Header().Add(
		httpheader.ContentType.String(),
		fmt.Sprintf("%s; charset=utf-8", mimetype.Json.String()),
	)
	_, err = w.Write([]byte("{ \"status\": \"success\"}"))
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
	w.Header().Add(
		httpheader.ContentType.String(),
		fmt.Sprintf("%s; charset=utf-8", mimetype.Json.String()),
	)
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
		app.apiNotFound(w, r)
		return
	}

	// convert to JSON bytes
	b, err := json.Marshal(purdoobahByName)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// send it out
	w.Header().Add(
		httpheader.ContentType.String(),
		fmt.Sprintf("%s; charset=utf-8", mimetype.Json.String()),
	)
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
		app.apiNotFound(w, r)
		return
	}

	// convert to JSON bytes
	b, err := json.Marshal(currentSection)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// send it out
	w.Header().Add(
		httpheader.ContentType.String(),
		fmt.Sprintf("%s; charset=utf-8", mimetype.Json.String()),
	)
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
		app.apiNotFound(w, r)
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
	w.Header().Add(
		httpheader.ContentType.String(),
		fmt.Sprintf("%s; charset=utf-8", mimetype.Json.String()),
	)
	_, err = w.Write(b)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) apiAllTraditions(w http.ResponseWriter, r *http.Request) {
	// get all traditions
	allTraditions, err := app.traditionService.All()
	if err != nil {
		app.serveError(w, err)
		return
	}

	// convert to JSON bytes
	b, err := json.Marshal(allTraditions)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// send it out
	w.Header().Add(
		httpheader.ContentType.String(),
		fmt.Sprintf("%s; charset=utf-8", mimetype.Json.String()),
	)
	_, err = w.Write(b)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *application) apiNotFound(w http.ResponseWriter, r *http.Request) {
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

package cachecontrol

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/purdoobahs/purdoobahs.com/internal/httpheader"
)

// Unix epoch time
var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	httpheader.Expires.String():      epoch,
	httpheader.CacheControl.String(): "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	httpheader.Pragma.String():       "no-cache",
	"X-Accel-Expires":                "0",
}

var etagHeaders = []string{
	httpheader.ETag.String(),
	httpheader.IfModifiedSince.String(),
	httpheader.IfMatch.String(),
	httpheader.IfNoneMatch.String(),
	httpheader.IfRange.String(),
	httpheader.IfUnmodifiedSince.String(),
}

type CacheControl struct {
	// Debug toggles debug log lines.
	Debug bool

	// ForeverCacheRoutePrefixes, if non-empty, tells the ForeverCache which route prefixes it is allowed to work on.
	ForeverCacheRoutePrefixes []string
}

func NewCacheControl() *CacheControl {
	return &CacheControl{
		Debug: false,

		// ForeverCache settings
		ForeverCacheRoutePrefixes: []string{},
	}
}

// ForeverCache is an HTTP server middleware which no-caches HTML, and caches every other MimeType for 1 year.
func (cc *CacheControl) ForeverCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if prefixes list is empty, ForeverCache all resources
		if len(cc.ForeverCacheRoutePrefixes) == 0 {
			cc.setForeverCache(w, r)
		} else {
			// since prefixes is not empty, cycle through all allowed prefixes to see if we have a match
			foundMatch := false
			for _, prefix := range cc.ForeverCacheRoutePrefixes {
				if strings.HasPrefix(r.URL.RequestURI(), prefix) {
					foundMatch = true
					break
				}
			}

			// if the requested URL doesn't match the allowed prefixes, set NoCache instead
			if foundMatch {
				cc.setForeverCache(w, r)
			} else {
				cc.setNoCache(w, r)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// NoCache is an HTTP server middleware which makes sure not a single resource is cached in any possible way.
func (cc *CacheControl) NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cc.setNoCache(w, r)
		next.ServeHTTP(w, r)
	})
}

func (cc *CacheControl) setForeverCache(w http.ResponseWriter, r *http.Request) {
	if cc.Debug {
		fmt.Printf("Caching forever: `%s`\n", r.URL.RequestURI())
	}

	// delete eTag headers
	cc.removeETagHeaders(w, r)

	// duration of 1 year
	foreverDuration := (24 * time.Hour) * 365

	// set Expires
	foreverUnix := time.Now().Add(foreverDuration).Format(time.RFC1123)
	w.Header().Set(httpheader.Expires.String(), foreverUnix)

	// set Cache-Control
	w.Header().Set(
		httpheader.CacheControl.String(),
		fmt.Sprintf(
			"public, max-age=%.0f, s-maxage=%.0f, must-revalidate, proxy-revalidate, immutable",
			foreverDuration.Seconds(),
			foreverDuration.Seconds(),
		),
	)
}

func (cc *CacheControl) setNoCache(w http.ResponseWriter, r *http.Request) {
	if cc.Debug {
		fmt.Printf("NoCache: `%s`\n", r.URL.RequestURI())
	}

	// delete all eTag related HTTP headers
	cc.removeETagHeaders(w, r)

	// set NoCache headers
	for httpHeaderKey, httpHeaderValue := range noCacheHeaders {
		w.Header().Set(httpHeaderKey, httpHeaderValue)
	}
}

func (cc *CacheControl) removeETagHeaders(w http.ResponseWriter, r *http.Request) {
	for _, eTagHttpHeader := range etagHeaders {
		if r.Header.Get(eTagHttpHeader) != "" {
			r.Header.Del(eTagHttpHeader)
		}
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/MagnusFrater/helmet"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info(fmt.Sprintf(
			"%s - %s %s %s",
			r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI()),
		)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serveError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func createHelmet() *helmet.Helmet {
	h := helmet.Empty()

	// sorted alphabetically
	h.ContentSecurityPolicy = helmet.NewContentSecurityPolicy(map[helmet.CSPDirective][]helmet.CSPSource{
		// fallback
		helmet.DirectiveDefaultSrc: {helmet.SourceNone},

		// these don't use 'default' as a fallback
		helmet.DirectiveBaseURI:    {helmet.SourceNone},
		helmet.DirectiveFormAction: {helmet.SourceNone},

		// these need to be not 'none'
		helmet.DirectiveFrameAncestors:      {helmet.SourceSelf},
		helmet.DirectiveFrameSrc:            {helmet.SourceSelf},
		helmet.DirectiveImgSrc:              {helmet.SourceSelf, helmet.SourceReportSample},
		helmet.DirectiveNavigateTo:          {helmet.SourceSelf, helmet.SourceReportSample},
		helmet.DirectiveObjectSrc:           {helmet.SourceSelf},
		helmet.DirectivePluginTypes:         {"application/pdf"},
		helmet.DirectiveReportTo:            {}, // TODO add support
		helmet.DeprecatedDirectiveReportURI: {"/csp-report"},
		helmet.DirectiveStyleSrc:            {helmet.SourceSelf, helmet.SourceReportSample},
		helmet.DirectiveScriptSrc:           {helmet.SourceSelf, helmet.SourceReportSample},
		helmet.DirectiveMediaSrc:            {helmet.SourceSelf},
		helmet.DirectiveConnectSrc:          {helmet.SourceSelf},

		// these are merely toggled and don't have options
		helmet.DirectiveBlockAllMixedContent: {},
	})

	h.XContentTypeOptions = helmet.XContentTypeOptionsNoSniff
	h.XDNSPrefetchControl = helmet.XDNSPrefetchControlOff
	h.XDownloadOptions = helmet.XDownloadOptionsNoOpen
	h.ExpectCT = helmet.NewExpectCT(30, true, "/expect-ct-report")

	h.FeaturePolicy = helmet.NewFeaturePolicy(map[helmet.FeaturePolicyDirective][]helmet.FeaturePolicyOrigin{
		helmet.DirectiveAccelerometer:     {helmet.OriginNone},
		helmet.DirectiveAutoplay:          {helmet.OriginNone},
		helmet.DirectiveCamera:            {helmet.OriginNone},
		helmet.DirectiveDocumentDomain:    {helmet.OriginNone},
		helmet.DirectiveEncryptedMedia:    {helmet.OriginNone},
		helmet.DirectiveFullscreen:        {helmet.OriginNone},
		helmet.DirectiveGeolocation:       {helmet.OriginNone},
		helmet.DirectiveGyroscope:         {helmet.OriginNone},
		helmet.DirectiveMagnetometer:      {helmet.OriginNone},
		helmet.DirectiveMicrophone:        {helmet.OriginNone},
		helmet.DirectiveMidi:              {helmet.OriginNone},
		helmet.DirectivePayment:           {helmet.OriginNone},
		helmet.DirectivePictureInPicture:  {helmet.OriginNone},
		helmet.DirectiveSyncXHR:           {helmet.OriginNone},
		helmet.DirectiveUSB:               {helmet.OriginNone},
		helmet.DirectiveXRSpacialTracking: {helmet.OriginNone},
	})

	h.XFrameOptions = helmet.XFrameOptionsDeny
	h.XPermittedCrossDomainPolicies = helmet.PermittedCrossDomainPoliciesNone
	h.XPoweredBy = helmet.NewXPoweredBy(true, "")
	h.ReferrerPolicy = helmet.NewReferrerPolicy(helmet.DirectiveStrictOriginWhenCrossOrigin)
	h.StrictTransportSecurity = helmet.NewStrictTransportSecurity(31536000, true, false)
	h.XXSSProtection = helmet.NewXXSSProtection(true, helmet.DirectiveModeBlock, "/xss-protection-report")

	return h
}

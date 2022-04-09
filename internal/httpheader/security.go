package httpheader

type Security HttpHeader

func (s Security) String() string {
	return string(s)
}

// List of security HTTP headers.
const (
	CrossOriginEmbedderPolicy       Security = "Cross-Origin-Embedder-Policy"
	CrossOriginOpenerPolicy         Security = "Cross-Origin-Opener-Policy"
	CrossOriginResourcePolicy       Security = "Cross-Origin-Resource-Policy"
	ContentSecurityPolicy           Security = "Content-Security-Policy"
	ContentSecurityPolicyReportOnly Security = "Content-Security-Policy-Report-Only"
	ExpectCt                        Security = "Expect-CT"
	FeaturePolicy                   Security = "Feature-Policy"
	StrictTransportSecurity         Security = "Strict-Transport-Security"
	UpgradeInsecureRequests         Security = "Upgrade-Insecure-Requests"
	XContentTypeOptions             Security = "X-Content-Type-Options"
	XDownloadOptions                Security = "X-Download-Options"
	XFrameOptions                   Security = "X-Frame-Options"
	XPermittedCrossDomainPolicies   Security = "X-Permitted-Cross-Domain-Policies"
	XPoweredBy                      Security = "X-Powered-By"
	XXssProtection                  Security = "X-XSS-Protection"
)

// List of experimental security HTTP headers.
const (
	ExperimentalOriginIsolation Security = "Origin-Isolation"
)

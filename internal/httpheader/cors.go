package httpheader

type Cors HttpHeader

func (c Cors) String() string {
	return string(c)
}

// List of CORS HTTP headers.
const (
	AccessControlAllowOrigin      Cors = "Access-Control-Allow-Origin"
	AccessControlAllowCredentials Cors = "Access-Control-Allow-Credentials"
	AccessControlAllowHeaders     Cors = "Access-Control-Allow-Headers"
	AccessControlAllowMethods     Cors = "Access-Control-Allow-Methods"
	AccessControlExposeHeaders    Cors = "Access-Control-Expose-Headers"
	AccessControlMaxAge           Cors = "Access-Control-Max-Age"
	AccessControlRequestHeaders   Cors = "Access-Control-Request-Headers"
	AccessControlRequestMethod    Cors = "Access-Control-Request-Method"
	Origin                        Cors = "Origin"
	TimingAllowOrigin             Cors = "Timing-Allow-Origin"
)

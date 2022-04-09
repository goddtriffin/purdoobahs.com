package httpheader

type UserAgentClientHints HttpHeader

func (uach UserAgentClientHints) String() string {
	return string(uach)
}

// List of experimental user agent client hint HTTP headers.
const (
	ExperimentalSecChUa                UserAgentClientHints = "Sec-CH-UA"
	ExperimentalSecChUaArch            UserAgentClientHints = "Sec-CH-UA-Arch"
	ExperimentalSecChUaBitness         UserAgentClientHints = "Sec-CH-UA-Bitness"
	ExperimentalSecChUaFullVersion     UserAgentClientHints = "Sec-CH-UA-Full-Version"
	ExperimentalSecChUaFullVersionList UserAgentClientHints = "Sec-CH-UA-Full-Version-List"
	ExperimentalSecChUaMobile          UserAgentClientHints = "Sec-CH-UA-Mobile"
	ExperimentalSecChUaModel           UserAgentClientHints = "Sec-CH-UA-Model"
	ExperimentalSecChUaPlatform        UserAgentClientHints = "Sec-CH-UA-Platform"
	ExperimentalSecChUaPlatformVersion UserAgentClientHints = "Sec-CH-UA-Platform-Version"
)

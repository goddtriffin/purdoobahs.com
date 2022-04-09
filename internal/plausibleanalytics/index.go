package plausibleanalytics

type Body struct {
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screen_width"`
}

func NewPlausibleAnalyticsBody(domain, name, url, referrer string, screenWidth int) *Body {
	return &Body{
		Domain:      domain,
		Name:        name,
		Url:         url,
		Referrer:    referrer,
		ScreenWidth: screenWidth,
	}
}

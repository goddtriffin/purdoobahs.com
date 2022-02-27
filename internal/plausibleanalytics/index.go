package plausibleanalytics

type Body struct {
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screen_width"`
}

func NewPlausibleAnalyticsBody(url, referrer string, screenWidth int) *Body {
	return &Body{
		Domain:      "purdoobahs.com",
		Name:        "pageview",
		Url:         url,
		Referrer:    referrer,
		ScreenWidth: screenWidth,
	}
}

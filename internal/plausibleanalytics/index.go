package plausibleanalytics

type PlausibleAnalyticsBody struct {
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screen_width"`
}

func NewPlausibleAnalyticsBody(url, referrer string, screen_width int) *PlausibleAnalyticsBody {
	return &PlausibleAnalyticsBody{
		Domain:      "purdoobahs.com",
		Name:        "pageview",
		Url:         url,
		Referrer:    referrer,
		ScreenWidth: screen_width,
	}
}

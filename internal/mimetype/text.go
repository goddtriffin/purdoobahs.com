package mimetype

type Text MimeType

func (t Text) String() string {
	return string(t)
}

const (
	Css            Text = "text/css"
	Csv            Text = "text/csv"
	Html           Text = "text/html"
	Calendar       Text = "text/calendar"
	JavascriptText Text = "text/javascript"
	Plain          Text = "text/plain"
	XmlText        Text = "text/xml"
)

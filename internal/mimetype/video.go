package mimetype

type Video MimeType

func (v Video) String() string {
	return string(v)
}

const (
	XMsVideo  Video = "video/x-msvideo"
	MpegVideo Video = "video/mpeg"
	Mp4       Video = "video/mp4"
	OggVideo  Video = "video/ogg"
	WebmVideo Video = "video/webm"
)

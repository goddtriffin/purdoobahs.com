package mimetype

type Audio MimeType

func (a Audio) String() string {
	return string(a)
}

const (
	Aac       Audio = "audio/aac"
	Midi      Audio = "audio/midi"
	XMidi     Audio = "audio/x-midi"
	MpegAudio Audio = "audio/mpeg"
	OggAudio  Audio = "audio/ogg"
	Opus      Audio = "audio/opus"
	Wav       Audio = "audio/wav"
	WebmAudio Audio = "audio/webm"
)

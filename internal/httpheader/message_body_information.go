package httpheader

type MessageBodyInformation HttpHeader

func (mbi MessageBodyInformation) String() string {
	return string(mbi)
}

// List of message body information HTTP headers.
const (
	ContentLength   MessageBodyInformation = "Content-Length"
	ContentType     MessageBodyInformation = "Content-Type"
	ContentEncoding MessageBodyInformation = "Content-Encoding"
	ContentLanguage MessageBodyInformation = "Content-Language"
	ContentLocation MessageBodyInformation = "Content-Location"
)

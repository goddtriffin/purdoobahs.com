package httpheader

type FetchMetadataRequestHeaders HttpHeader

func (fmrh FetchMetadataRequestHeaders) String() string {
	return string(fmrh)
}

// List of fetch metadata request headers HTTP headers.
const (
	SecFetchSite                   FetchMetadataRequestHeaders = "Sec-Fetch-Site"
	SecFetchMode                   FetchMetadataRequestHeaders = "Sec-Fetch-Mode"
	SecFetchUser                   FetchMetadataRequestHeaders = "Sec-Fetch-User"
	SecFetchDest                   FetchMetadataRequestHeaders = "Sec-Fetch-Dest"
	ServiceWorkerNavigationPreload FetchMetadataRequestHeaders = "Service-Worker-Navigation-Preload"
)

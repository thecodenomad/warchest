package query

import (
	"net/http"
)

// CBBaseURL is the baseurl for all coinbase API calls
const CBBaseURL = "https://api.coinbase.com"

var (
	// ErrDecoding occurs when failing to decode an api response
	ErrDecoding = Error("failed decoding response")

	// ErrOnUnmarshall occurs when failing to build an object from an api response
	ErrOnUnmarshall = Error("failed to unmarshall")

	// ErrConnection occurs when a request is interupted or fails
	ErrConnection = Error("error during request")
)

// Error is the helper method that produces the errors above
func (e Error) Error() string {
	return string(e)
}

// Error the object for query errors
type Error string

// HTTPClient interface is an internal interface useful for testability
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

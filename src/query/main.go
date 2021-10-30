package query

import (
	"net/http"
)

const CBBaseURL = "https://api.coinbase.com"

var (
	ErrDecoding     = QueryError("failed decoding response")
	ErrOnUnmarshall = QueryError("failed to unmarshall")
	ErrConnection   = QueryError("error during request")
)

func (q QueryError) Error() string {
	return string(q)
}

type QueryError string

// HttpClient interface is an internal interface useful for testability
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

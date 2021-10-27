package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

var (
	cbaccess_key       = "CB-ACCESS-KEY"
	cbaccess_sign      = "CB-ACCESS-SIGN"
	cbaccess_timestamp = "CB-ACCESS-TIMESTAMP"
)

type CBAuth struct {
	// private!
	apiKey        string
	apiSecret     string
	requestMethod string
	requestBody   string
	requestPath   string
	timestamp     int
}

// Ref: https://developers.coinbase.com/docs/wallet/api-key-authentication
// All REST requests must contain the following headers:
//
// CB-ACCESS-KEY API key as a string
// CB-ACCESS-SIGN Message signature (see below)
// CB-ACCESS-TIMESTAMP Timestamp for your request
//
// curl https://api.coinbase.com/v2/user \
// --header "CB-ACCESS-KEY: <your api key>" \
// --header "CB-ACCESS-SIGN: <the user generated message signature>" \
// --header "CB-ACCESS-TIMESTAMP: <a timestamp for your request>"
//
// The CB-ACCESS-SIGN header is generated by creating a sha256 HMAC using the secret key
// on the prehash string timestamp + method + requestPath + body (where + represents string
// concatenation). The timestamp value is the same as the CB-ACCESS-TIMESTAMP header.
//
// The body is the request body string. It is omitted if there is no request body (typically
// for GET requests).
//
// The method should be UPPER CASE.
//
// The requestPath is the full path and query parameters of the URL, e.g.:
//   /v2/exchange-rates?currency=USD.
//
// The CB-ACCESS-TIMESTAMP header MUST be number of seconds since Unix Epoch in UTC.
//
// Your timestamp must be within 30 seconds of the API service time, or your request will be
// considered expired and rejected. We recommend using the time API endpoint to query for
// the API server time if you believe there may be a time skew between your server and the API
// servers.

func (c *CBAuth) NewAuthMap() map[string]string {

	// Setup return object
	headers := map[string]string{}

	// Generate new timestamp for this call
	timestamp := int(time.Now().Unix())

	// Sign at the dotted line...
	sigText := "%d" + c.requestMethod + c.requestPath + c.requestBody
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(sigText))
	signature := hex.EncodeToString(h.Sum(nil))

	// Apply the things
	headers[cbaccess_key] = c.apiKey
	headers[cbaccess_timestamp] = fmt.Sprintf("%d", timestamp)
	headers[cbaccess_sign] = signature

	return headers
}
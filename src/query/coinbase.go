package query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const CBBaseURL = "https://api.coinbase.com"
const CBExchangeRateUrl = "/v2/exchange-rates"
const CBUserUrl = "/v2/user"
const CBTransactionUrl = "/v2/accounts/:account_id/transactions"

var (
	ErrDecoding     = QueryError("failed decoding response")
	ErrOnUnmarshall = QueryError("failed to unmarshall")
	ErrConnection   = QueryError("error during request")
)

func (q QueryError) Error() string {
	return string(q)
}

type QueryError string

type CoinInfoResponse struct {
	Data CoinInfo `json:"data"`
}

type CoinInfo struct {
	Currency      string `json:"currency"`
	ExchangeRates Rates  `json:"rates"`
}

type Rates struct {
	EUR float64 `json:"EUR,string"`
	GBP float64 `json:"GBP,string"`
	USD float64 `json:"USD,string"`
}

type CBUserResponse struct {
	Data struct {
		Id              string      `json:"id"`
		Name            string      `json:"name"`
		Username        string      `json:"username"`
		ProfileLocation interface{} `json:"profile_location"`
		ProfileBio      interface{} `json:"profile_bio"`
		ProfileUrl      string      `json:"profile_url"`
		AvatarUrl       string      `json:"avatar_url"`
		Resource        string      `json:"resource"`
		ResourcePath    string      `json:"resource_path"`
	} `json:"data"`
}

// HttpClient interface is an internal interface useful for testability
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CBRetrieveCoinData will return exchange rates for a given Crypto Currency Symbol
func CBRetrieveCoinData(symbol string, client HttpClient) (CoinInfo, error) {
	url := CBBaseURL + CBExchangeRateUrl + "?currency=" + symbol

	req, err := http.NewRequest("GET", url, nil)

	// Retrieve response
	resp, err := client.Do(req)

	//TODO: Create custom error for connectivity failures
	if err != nil {
		log.Printf("Hit error on retrieval: %s", err)
		return CoinInfo{}, ErrConnection
	}
	defer resp.Body.Close()

	cResp := CoinInfoResponse{}

	//TODO: Create custom error for failures to read respBody as a string
	bodyAsStr, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body of error: %s", err)
		return CoinInfo{}, ErrDecoding
	}

	//TODO: Create custom error for unmarshalling issues
	if err := json.Unmarshal([]byte(bodyAsStr), &cResp); err != nil {
		log.Printf("%s", err)
		return CoinInfo{}, ErrOnUnmarshall
	}

	return cResp.Data, err
}

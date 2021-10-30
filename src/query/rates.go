package query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const CBExchangeRateUrl = "/v2/exchange-rates"

type CoinInfoResp struct {
	Info CoinInfo `json:"data"`
}

type CoinInfo struct {
	Currency string    `json:"currency"`
	Rates    CoinRates `json:"rates"`
}

type CoinRates struct {
	EUR float64 `json:"EUR,string"`
	GBP float64 `json:"GBP,string"`
	USD float64 `json:"USD,string"`
}

// CBRetrieveCoinRate will return exchange rates for a given Crypto Currency Symbol
func CBRetrieveCoinRate(symbol string, client HttpClient) (CoinRates, error) {
	url := CBBaseURL + CBExchangeRateUrl + "?currency=" + symbol

	req, err := http.NewRequest("GET", url, nil)

	// Retrieve response
	resp, err := client.Do(req)

	//TODO: Create custom error for connectivity failures
	if err != nil {
		log.Printf("Hit error on retrieval: %s", err)
		return CoinRates{}, ErrConnection
	}
	defer resp.Body.Close()

	cResp := CoinInfoResp{}

	//TODO: Create custom error for failures to read respBody as a string
	bodyAsStr, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body of error: %s", err)
		return CoinRates{}, ErrDecoding
	}

	//TODO: Create custom error for unmarshalling issues
	if err := json.Unmarshal([]byte(bodyAsStr), &cResp); err != nil {
		log.Printf("%s", err)
		return CoinRates{}, ErrOnUnmarshall
	}

	return cResp.Info.Rates, err
}

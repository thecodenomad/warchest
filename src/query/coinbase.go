package query

import (
	"encoding/json"
	"log"
	"net/http"
)

const ExchangeRateUrl = "https://api.coinbase.com/v2/exchange-rates"

type JSONResponse struct {
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

//type Wallet map[string]map[string]float64

// CBRetrieveCoinData will return exchange rates for a given Crypto Currency Symbol
// TODO: Change the structure so that this is a struct method
func CBRetrieveCoinData(symbol string) (CoinInfo, error) {

	url := ExchangeRateUrl + "?currency=" + symbol
	resp, err := http.Get(url)

	//TODO: Create custom error for failure to decode
	if err != nil {
		return CoinInfo{}, err
	}
	defer resp.Body.Close()

	var cResp JSONResponse

	//TODO: Create custom error for failure to decode
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Printf("error: %s", err)
		return CoinInfo{}, err
	}

	return cResp.Data, err
}

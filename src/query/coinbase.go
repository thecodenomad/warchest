package query

import (
	"encoding/json"
	"log"
	"net/http"
)

const exchangeRateUrl = "https://api.coinbase.com/v2/exchange-rates"

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

// RetrieveCoinData will return exchange rates for a given Crypto Curraency Symbol
func RetrieveCoinData(symbol string) (CoinInfo, error) {

	url := exchangeRateUrl + "?currency=" + symbol
	resp, err := http.Get(url)

	// TODO: Create custom error for failed response
	if err != nil {
		log.Fatal("ooopsss an error occurred, please try again")
	}
	defer resp.Body.Close()

	var cResp JSONResponse

	//TODO: Create custom error for failure to decode
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Printf("error: %s", err)
		log.Fatal("ooopsss! an error occurred, please try again")
		return CoinInfo{}, err
	}

	return cResp.Data, err
}

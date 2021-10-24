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
	EURO string `json:"EUR"`
	GBP  string `json:"GBP"`
	USD  string `json:"USD"`
}

// RetrieveCoinData will return exchange rates for a given Crypto Curraency Symbol
func RetrieveCoinData(symbol string) CoinInfo {

	url := exchangeRateUrl + "?currency=" + symbol
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("ooopsss an error occurred, please try again")
	}
	defer resp.Body.Close()

	var cResp JSONResponse

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Printf("error: %s", err)
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	return cResp.Data
}

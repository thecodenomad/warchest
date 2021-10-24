package query

import (
	"encoding/json"
	"log"
	"net/http"
	"warchest/src/config"
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

type Wallet struct {
	Coins     []Coin
	NetProfit float64
}

// Coin used to store an individual coins profit
// TODO: is there a better way to keep this DRY? ref PurchasedCoins
type Coin struct {
	CoinSymbol     string
	CurrentRateUSD float64
	CurrentRateGBP float64
	CurrentRateEUR float64
	Transactions   []CoinTransaction
}

type CoinTransaction struct {
	NumCoins       float64
	PurchasedPrice float64
	TransactionFee float64
}

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

// CalculateCoinProfit calculates the net profit across all coin purchases
func CalculateCoinProfit(coin Coin, info CoinInfo) (float64, error) {

	totalNumCoins := 0.0
	totalExpense := 0.0

	// Need to iterate over each coin purchase and add together to find
	// value
	//  -- Amount*PurchasedPrice + transaction fees
	for _, transaction := range coin.Transactions {
		totalNumCoins += transaction.NumCoins
		totalExpense += transaction.NumCoins*transaction.PurchasedPrice + transaction.TransactionFee
	}

	// TODO: extend to support other currencies (crypto and national)
	totalProfit := info.ExchangeRates.USD*totalNumCoins - totalExpense

	return totalProfit, nil
}

func CalculateNetProfit(coins config.PurchasedCoins) Wallet {

	// Dummy Data for compilation
	transactions := []CoinTransaction{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
	}
	testCoin := Coin{"ETH", 30.0, 40.0, 50.0, transactions}

	return Wallet{[]Coin{testCoin}, 100.0}
}

package wallet

import (
	"fmt"
	"warchest/src/query"
)

type Wallet struct {
	Coins     []Coin
	NetProfit float64
}

// Coin used to store an individual coins profit
// TODO: is there a better way to keep this DRY? ref PurchasedCoins
type Coin struct {
	CoinSymbol       string
	CurrentRateUSD   float64
	CurrentRateGBP   float64
	CurrentRateEUR   float64
	CurrentNetProfit float64
	Transactions     []CoinTransaction
}

type CoinTransaction struct {
	NumCoins       float64
	PurchasedPrice float64
	TransactionFee float64
}

func (c *Coin) UpdateRates() {

	coinInfo, err := query.RetrieveCoinData(c.CoinSymbol)
	if err != nil {
		fmt.Printf("Failed to retrieve market rates for %s\n", c.CoinSymbol)
		return
	}

	// Update the rates
	c.CurrentRateEUR = coinInfo.ExchangeRates.EUR
	c.CurrentRateGBP = coinInfo.ExchangeRates.GBP
	c.CurrentRateUSD = coinInfo.ExchangeRates.USD
}

// CalculateCoinProfit calculates the net profit across all coin purchases (assume coins already have rates updated)
func CalculateCoinProfit(coin Coin) (float64, error) {

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
	totalProfit := coin.CurrentRateUSD*totalNumCoins - totalExpense

	return totalProfit, nil
}

func CalculateNetProfit(wallet Wallet) (float64, error) {

	netProfit := 0.0

	for _, coin := range wallet.Coins {
		coinProfit, err := CalculateCoinProfit(coin)
		if err != nil {
			fmt.Printf("Failed to calculate profit for %s\n", coin.CoinSymbol)
			return 0.0, err
		}
		netProfit += coinProfit
	}

	return netProfit, nil
}

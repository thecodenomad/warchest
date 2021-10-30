package wallet

import (
	"fmt"
	"log"
	"warchest/src/query"
)

type Wallet struct {
	Coins     []WarchestCoin
	NetProfit float64
}

// Coin used to store an individual coins profit
// TODO: is there a better way to keep this DRY? ref PurchasedCoins
type WarchestCoin struct {
	Cost         float64
	Amount       float64
	Profit       float64
	Rates        query.CoinRates
	CoinSymbol   string
	Transactions []CoinTransaction
}

type CoinTransaction struct {
	NumCoins       float64
	PurchasedPrice float64
	TransactionFee float64
}

//UpdateRates updates a coin's current exchange rate
func (w *WarchestCoin) UpdateRates(client query.HttpClient) {

	coinRates, err := query.CBRetrieveCoinRate(w.CoinSymbol, client)
	if err != nil {
		log.Printf("Failed to retrieve market rates for %s\n", w.CoinSymbol)
		// Reset instead of erroring
		w.Rates.EUR = 0.0
		w.Rates.GBP = 0.0
		w.Rates.USD = 0.0
		return
	}

	// Update the rates
	w.Rates.EUR = coinRates.EUR
	w.Rates.GBP = coinRates.GBP
	w.Rates.USD = coinRates.USD
}

//UpdateCost updates a coin's initial purchase cost from the coins transactions
func (w *WarchestCoin) UpdateCost() {
	totalNumCoins := 0.0
	totalExpense := 0.0

	log.Printf("There are %d %s transactions in your wallet, calculating...\n", len(w.Transactions), w.CoinSymbol)

	for _, transaction := range w.Transactions {
		totalNumCoins += transaction.NumCoins
		totalExpense += transaction.NumCoins*transaction.PurchasedPrice + transaction.TransactionFee
	}

	w.Amount = totalNumCoins
	w.Cost = totalExpense
}

//UpdateProfit updates a coin's net profit value
func (w *WarchestCoin) UpdateProfit() {
	currentValue := w.Rates.USD*w.Amount - w.Cost
	w.Profit = currentValue
}

//Update runs all internal updates to get the latest value of a particular coin in a wallet
func (w *WarchestCoin) Update(client query.HttpClient) {
	w.UpdateCost()
	w.UpdateRates(client)
	w.UpdateProfit()
}

//Banner prints out a stats banner for the coin
func (w *WarchestCoin) Banner() {
	fmt.Printf("\tCurrent rate for %s: %.6f\n", w.CoinSymbol, w.Rates.USD)
	fmt.Printf("\tInitial Cost of %s: %.6f\n", w.CoinSymbol, w.Cost)
	fmt.Printf("\tTotal Amount of %s: %.6f\n", w.CoinSymbol, w.Amount)
	fmt.Printf("\tCurrent cost of %s: %.6f\n", w.CoinSymbol, w.Amount*w.Rates.USD)
	fmt.Printf("\tTotal profit for %s: %.6f\n", w.CoinSymbol, w.Profit)
}

// CalculateNetProfit will calculate the total profit for the coins in the provided Wallet
func CalculateNetProfit(wallet Wallet, client query.HttpClient) (float64, error) {
	netProfit := 0.0

	log.Printf("There are %d coin(s) in your wallet, calculating...\n", len(wallet.Coins))
	for _, coin := range wallet.Coins {
		// Make sure we have the latest rates
		coin.Update(client)

		// Present stats for coin
		coin.Banner()
		netProfit += coin.Profit
	}
	return netProfit, nil
}

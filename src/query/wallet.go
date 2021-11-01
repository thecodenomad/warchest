package query

import (
	"fmt"
	"log"
	"warchest/src/auth"
)

//
// Warchest Objects
//////////////////////

// Wallet is the main object consumed by warchest that includes all coins and their transactions
type Wallet struct {
	Coins     []WarchestCoin
	NetProfit float64
}

// WarchestCoin a coin object that includes stats and transactions for purchased coins
// TODO: is there a better way to keep this DRY? ref PurchasedCoins
type WarchestCoin struct {
	AccountID    string
	Cost         float64
	Amount       float64
	Profit       float64
	Rates        CoinRates
	CoinSymbol   string
	Transactions []CoinTransaction
}

// CoinTransaction is an individual transaction made for a given type of coin
type CoinTransaction struct {
	NumCoins       float64
	PurchasedPrice float64
	TransactionFee float64
}

// UpdateTransactions method will retrieve the transactions for a given coin
func (w *WarchestCoin) UpdateTransactions(cbAuth auth.CBAuth, client HTTPClient) {

	transactions, err := CBCoinTransactions(w.AccountID, cbAuth, client)
	if err != nil {
		fmt.Printf("Failed retreiving transactions: %s", err)
		w.Transactions = []CoinTransaction{}
	}

	coinTransactions := []CoinTransaction{}

	fmt.Printf("There are %d transactions for %s\n", len(transactions), w.CoinSymbol)

	for _, cbTransaction := range transactions {
		fmt.Printf("Adding transaction for %s\n", cbTransaction.Amount.Currency)
		coinTransactions = append(coinTransactions, cbTransaction.ToCoinTransaction())
	}
	w.Transactions = coinTransactions
}

//UpdateRates updates a coin's current exchange rate
func (w *WarchestCoin) UpdateRates(client HTTPClient) {

	coinRates, err := CBRetrieveCoinRates(w.CoinSymbol, client)
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
	fmt.Printf("USD Rate: %.6f\n", w.Rates.USD)
	fmt.Printf("Amount: %.6f\n", w.Amount)
	fmt.Printf("Cost: %.6f\n", w.Cost)
	currentValue := w.Rates.USD*w.Amount - w.Cost
	w.Profit = currentValue
}

//Update runs all internal updates to get the latest value of a particular coin in a wallet
func (w *WarchestCoin) Update(cbAuth auth.CBAuth, client HTTPClient) {
	w.UpdateTransactions(cbAuth, client)
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
func CalculateNetProfit(wallet Wallet, cbAuth auth.CBAuth, client HTTPClient) (float64, error) {
	netProfit := 0.0

	log.Printf("There are %d coin(s) in your wallet, calculating...\n", len(wallet.Coins))
	for _, coin := range wallet.Coins {
		// Make sure we have the latest rates
		coin.Update(cbAuth, client)

		// Present stats for coin
		coin.Banner()
		netProfit += coin.Profit
	}
	return netProfit, nil
}

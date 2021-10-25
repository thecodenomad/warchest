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
	Amount           float64
	Cost             float64
	Profit           float64
	Transactions     []CoinTransaction
}

type CoinTransaction struct {
	NumCoins       float64
	PurchasedPrice float64
	TransactionFee float64
}

//UpdateRates updates a coin's current exchange rate
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

//UpdateCost updates a coin's initial purchase cost from the coins transactions
func (c *Coin) UpdateCost() {
	totalNumCoins := 0.0
	totalExpense := 0.0

	fmt.Printf("There are %d %s transactions in your wallet, calculating...\n", len(c.Transactions), c.CoinSymbol)

	for _, transaction := range c.Transactions {
		totalNumCoins += transaction.NumCoins
		totalExpense += transaction.NumCoins*transaction.PurchasedPrice + transaction.TransactionFee

		//fmt.Printf("Transaction has %.6f of %s\n", transaction.NumCoins, c.CoinSymbol)
		//fmt.Printf("Transaction has %.6f cost\n", totalExpense)
	}

	c.Amount = totalNumCoins
	c.Cost = totalExpense
}

//UpdateProfit updates a coin's net profit value
func (c *Coin) UpdateProfit() {
	currentValue := c.CurrentRateUSD*c.Amount - c.Cost
	c.Profit = currentValue
}

//Update runs all internal updates to get the latest value of a particular coin in a wallet
func (c *Coin) Update() {
	c.UpdateCost()
	c.UpdateRates()
	c.UpdateProfit()
}

//Banner prints out a stats banner for the coin
func (c *Coin) Banner() {
	fmt.Printf("\tCurrent rate for %s: %.6f\n", c.CoinSymbol, c.CurrentRateUSD)
	fmt.Printf("\tInitial Cost of %s: %.6f\n", c.CoinSymbol, c.Cost)
	fmt.Printf("\tTotal Amount of %s: %.6f\n", c.CoinSymbol, c.Amount)
	fmt.Printf("\tCurrent cost of %s: %.6f\n", c.CoinSymbol, c.Amount*c.CurrentRateUSD)
	fmt.Printf("\tTotal profit for %s: %.4f\n", c.CoinSymbol, c.Profit)
}

// CalculateNetProfit will calculate the total profit for the coins in the provided Wallet
func CalculateNetProfit(wallet Wallet) (float64, error) {
	netProfit := 0.0

	fmt.Printf("There are %d coin(s) in your wallet, calculating...\n", len(wallet.Coins))
	for _, coin := range wallet.Coins {
		// Make sure we have the latest rates
		coin.Update()

		// Present stats for coin
		coin.Banner()
		netProfit += coin.Profit
	}
	return netProfit, nil
}

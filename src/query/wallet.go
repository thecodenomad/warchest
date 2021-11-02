package query

import (
	"log"
	"warchest/src/auth"
)

//
// Warchest Objects
//////////////////////

// Wallet is the main object consumed by warchest that includes all coins and their transactions
type Wallet struct {
	Coins     map[string]WarchestCoin `json:"coins"`
	NetProfit float64                 `json:"net_profit"`
}

// WarchestCoin a coin object that includes stats and transactions for purchased coins
// TODO: is there a better way to keep this DRY? ref PurchasedCoins
type WarchestCoin struct {
	AccountID    string            `json:"account_id"`
	Cost         float64           `json:"cost"`
	Amount       float64           `json:"amount"`
	Profit       float64           `json:"profit"`
	Rates        CoinRates         `json:"rates"`
	Symbol       string            `json:"symbol""`
	Transactions []CoinTransaction `json:"transactions"`
}

// CoinTransaction is an individual transaction made for a given type of coin
type CoinTransaction struct {
	NumCoins       float64 `json:"num_coins"`
	PurchasedPrice float64 `json:"purchased_price"`
	TransactionFee float64 `json:"transaction_fee"`
}

// getSupportedCoins is an internal helper function that returns the currently supported coins for warchest
// TODO: this list will expand, this is just a way of bypassing coins that may have interest accruing
func getSupportedCoins() []string {
	return []string{"DOGE", "SHIB"}
}

// IsSupportedCoin is a helper method to determine if a coin is supported or not
func (w *WarchestCoin) IsSupportedCoin() bool {
	supportedCoins := getSupportedCoins()
	for _, supportedCoin := range supportedCoins {
		if w.Symbol == supportedCoin {
			return true
		}
	}
	return false
}

// UpdateTransactions method will retrieve the transactions for a given coin
func (w *WarchestCoin) UpdateTransactions(cbAuth auth.CBAuth, client HTTPClient) {

	transactions, err := CBCoinTransactions(w.AccountID, cbAuth, client)
	if err != nil {
		log.Printf("Failed retreiving transactions: %s", err)
		w.Transactions = []CoinTransaction{}
	}

	coinTransactions := []CoinTransaction{}

	log.Printf("There are %d transactions for %s\n", len(transactions), w.Symbol)

	for _, cbTransaction := range transactions {
		log.Printf("Adding transaction for %s\n", cbTransaction.Amount.Currency)
		log.Printf("NumCoins: %.14f\n", cbTransaction.Amount.Amount)
		log.Printf("PurchasedPrices: %.14f\n", cbTransaction.NativeAmount.Amount)
		coinTransactions = append(coinTransactions, cbTransaction.ToCoinTransaction())
	}
	w.Transactions = coinTransactions
}

//UpdateRates updates a coin's current exchange rate
func (w *WarchestCoin) UpdateRates(client HTTPClient) {

	coinRates, err := CBRetrieveCoinRates(w.Symbol, client)
	if err != nil {
		log.Printf("Failed to retrieve market rates for %s\n", w.Symbol)
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

	for _, transaction := range w.Transactions {
		totalNumCoins += transaction.NumCoins
		// TODO: need 'purchased rate' instead
		// NOTE: CB API - fee is in the total price
		totalExpense += transaction.PurchasedPrice
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
func (w *WarchestCoin) Update(cbAuth auth.CBAuth, client HTTPClient) {
	w.UpdateTransactions(cbAuth, client)
	w.UpdateCost()
	w.UpdateRates(client)
	w.UpdateProfit()
}

//Banner prints out a stats banner for the coin
func (w *WarchestCoin) Banner() {
	log.Printf("\tCurrent rate for %s: %.14f\n", w.Symbol, w.Rates.USD)
	log.Printf("\tInitial Cost of %s: %.14f\n", w.Symbol, w.Cost)
	log.Printf("\tTotal Amount of %s: %.14f\n", w.Symbol, w.Amount)
	log.Printf("\tCurrent cost of %s: %.14f\n", w.Symbol, w.Amount*w.Rates.USD)
	log.Printf("\tTotal profit for %s: %.14f\n", w.Symbol, w.Profit)
}

// UpdateNetProfit will calculate the total profit for the coins in the provided Wallet
func (w *Wallet) UpdateNetProfit(cbAuth auth.CBAuth, client HTTPClient) (float64, error) {
	netProfit := 0.0

	log.Printf("There are %d coin(s) in your wallet, calculating...\n", len(w.Coins))
	for _, coin := range w.Coins {

		// If there aren't transactions for this coin, retrieve them
		if len(coin.Transactions) < 1 {
			coin.UpdateTransactions(cbAuth, client)
		}

		log.Printf("Updating Cost, Current Rates, and Profit for %s", coin.Symbol)
		// Make sure cost is calculated
		coin.UpdateCost()

		// Make sure we have the latest rates
		coin.UpdateRates(client)

		// Now recalculate based on the updated rates
		coin.UpdateProfit()

		//coin.Update(cbAuth, client)

		// Present stats for coin
		// coin.Banner()
		netProfit += coin.Profit
	}

	// Make sure objects value is updated
	w.NetProfit = netProfit
	return netProfit, nil
}

// UpdateCoinRates will update the rates for all coins in a given wallet
func (w *Wallet) UpdateCoinRates(client HTTPClient) {
	for _, coin := range w.Coins {
		coin.UpdateRates(client)
	}
}

// GetWarchestCoins will retrieve all of the 'accounts' and convert them into a map of WarchestCoins
func GetWarchestCoins(cbAuth auth.CBAuth, client HTTPClient) (map[string]WarchestCoin, error) {

	accountResp, err := CBRetrieveAccounts(cbAuth, client)
	if err != nil {
		log.Printf("Failed to retrieve accounts: %s", err)
		return map[string]WarchestCoin{}, err
	}

	log.Printf("There are %d accounts to look through", len(accountResp.Accounts))

	coins := map[string]WarchestCoin{}
	for _, account := range accountResp.Accounts {
		coinToAdd := WarchestCoin{
			AccountID:    account.ID,
			Cost:         0.0,
			Amount:       0.0,
			Profit:       0.0,
			Rates:        CoinRates{},
			Symbol:       account.Currency.Code,
			Transactions: []CoinTransaction{},
		}

		// Check if this is a supported coin first
		if !coinToAdd.IsSupportedCoin() {
			log.Printf("Skipping unsupported coin: %s\n", coinToAdd.Symbol)
			continue
		}

		// Make sure coin updates appropriately
		// TODO: Add error handling around this as update _could_ fail
		coinToAdd.Update(cbAuth, client)

		// Add to the map!
		log.Printf("Adding coin %s to the list of coins", coinToAdd.Symbol)
		coins[account.Currency.Code] = coinToAdd
	}

	return coins, nil
}

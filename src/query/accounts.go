package query

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"warchest/src/auth"
)

// CBAccountsURL is the path to the GET accounts API call
const CBAccountsURL = "/v2/accounts"

// CBAccountsResp is the response object returned by retreiving all accounts
// TODO: need to handle pagination if there are that many accounts
type CBAccountsResp struct {
	Pagination struct {
		EndingBefore         interface{} `json:"ending_before"`
		StartingAfter        interface{} `json:"starting_after"`
		PreviousEndingBefore interface{} `json:"previous_ending_before"`
		NextStartingAfter    string      `json:"next_starting_after"`
		Limit                int         `json:"limit"`
		Order                string      `json:"order"`
		PreviousURI          interface{} `json:"previous_uri"`
		NextURI              string      `json:"next_uri"`
	} `json:"pagination"`
	Accounts []CBAccount `json:"data"`
}

// CBAccount is an individual account object
type CBAccount struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Primary  bool   `json:"primary"`
	Type     string `json:"type"`
	Currency struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Color        string `json:"color"`
		SortIndex    int    `json:"sort_index"`
		Exponent     int    `json:"exponent"`
		Type         string `json:"type"`
		AddressRegex string `json:"address_regex"`
		AssetID      string `json:"asset_id"`
		Slug         string `json:"slug"`
	} `json:"currency"`
	Balance struct {
		Amount   float64 `json:"amount,string"`
		Currency string  `json:"currency"`
	} `json:"balance"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Resource         string    `json:"resource"`
	ResourcePath     string    `json:"resource_path"`
	AllowDeposits    bool      `json:"allow_deposits"`
	AllowWithdrawals bool      `json:"allow_withdrawals"`
}

// CBRetrieveAccounts is a function that will retrieve all accounts associated with the account
func CBRetrieveAccounts(cbAuth auth.CBAuth, client HTTPClient) (CBAccountsResp, error) {

	url := CBBaseURL + CBAccountsURL
	authHeaders := cbAuth.NewAuthMap("GET", "", CBAccountsURL)
	req, err := http.NewRequest("GET", url, nil)

	// Set auth headers
	for key, value := range authHeaders {
		log.Printf("Setting %s to %s", key, value)
		req.Header.Add(key, value)
	}

	// per https://developers.coinbase.com/api/v2?shell#versioning
	t := time.Now()
	utcDate := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	req.Header.Add("CB-VERSION", utcDate)

	// Retrieve response
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Hit error on retrieval: %s", err)
		return CBAccountsResp{}, ErrConnection
	}
	defer resp.Body.Close()

	cResp := CBAccountsResp{}

	bodyAsStr, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body of error: %s", err)
		return CBAccountsResp{}, ErrDecoding
	}

	if err := json.Unmarshal([]byte(bodyAsStr), &cResp); err != nil {
		log.Printf("Body as string: \n%s\n", bodyAsStr)
		log.Printf("Couldn't unmarshall: %s", err)
		return CBAccountsResp{}, ErrOnUnmarshall
	}
	//log.Printf("Body as string: \n%s\n", bodyAsStr)

	return cResp, err
}

////UpdateRates updates a coin's current exchange rate
//func (c *CBAccount) UpdateRates(client HTTPClient) {
//
//	coinInfo, err := CBRetrieveCoinRate(c.Currency.Code, client)
//	if err != nil {
//		log.Printf("Failed to retrieve market rates for %s\n", c.Name)
//		// Reset instead of erroring
//		c.CurrentRateUSD = 0.0
//		return
//	}
//
//	c.CurrentRateUSD = coinInfo.ExchangeRates.USD
//}
//
////UpdateCost updates a coin's initial purchase cost from the coins transactions
//func (c *CBAccount) UpdateCost() {
//	totalNumCoins := 0.0
//	totalExpense := 0.0
//
//	log.Printf("There are %d %s transactions in your wallet, calculating...\n", len(c.Transactions), c.Name)
//
//	for _, transaction := range c.Transactions {
//		totalNumCoins += transaction.NumCoins
//		totalExpense += transaction.NumCoins*transaction.PurchasedPrice + transaction.TransactionFee
//	}
//
//	c.Amount = totalNumCoins
//	c.Cost = totalExpense
//}
//
////UpdateProfit updates a coin's net profit value
//func (c *Coin) UpdateProfit() {
//	currentValue := c.CurrentRateUSD*c.Amount - c.Cost
//	c.Profit = currentValue
//}
//
////Update runs all internal updates to get the latest value of a particular coin in a wallet
//func (c *Coin) Update(client query.HTTPClient) {
//	c.UpdateCost()
//	c.UpdateRates(client)
//	c.UpdateProfit()
//}
//
////Banner prints out a stats banner for the coin
//func (c *Coin) Banner() {
//	fmt.Printf("\tCurrent rate for %s: %.6f\n", c.CoinSymbol, c.CurrentRateUSD)
//	fmt.Printf("\tInitial Cost of %s: %.6f\n", c.CoinSymbol, c.Cost)
//	fmt.Printf("\tTotal Amount of %s: %.6f\n", c.CoinSymbol, c.Amount)
//	fmt.Printf("\tCurrent cost of %s: %.6f\n", c.CoinSymbol, c.Amount*c.CurrentRateUSD)
//	fmt.Printf("\tTotal profit for %s: %.6f\n", c.CoinSymbol, c.Profit)
//}
//
//// CalculateNetProfit will calculate the total profit for the coins in the provided Wallet
//func CalculateNetProfit(wallet Wallet, client query.HTTPClient) (float64, error) {
//	netProfit := 0.0
//
//	log.Printf("There are %d coin(s) in your wallet, calculating...\n", len(wallet.Coins))
//	for _, coin := range wallet.Coins {
//		// Make sure we have the latest rates
//		coin.Update(client)
//
//		// Present stats for coin
//		coin.Banner()
//		netProfit += coin.Profit
//	}
//	return netProfit, nil
//}

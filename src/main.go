package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
	"warchest/src/auth"
	// "warchest/src/config"
	"warchest/src/query"
)

// FailedLoadConfigRC Return code for failing to load the Warchest configuration
const FailedLoadConfigRC = 2

// FailedRetrievingData Return code for failing the retrieval of coin data
const FailedRetrievingData = 3

// FailedCalculatingWallet Return code for failing caculation of wallet
const FailedCalculatingWallet = 3

// WarchestConfigEnv is the environment variable that will point to coin transactions used by Warchest
const WarchestConfigEnv = "WARCHEST_CONFIG"

// CbAPIKey is the api key established via your coinbase profile
const CbAPIKey = "CB_API_KEY"

// CbAPISecret is the secret associated with the cbAPIKey
const CbAPISecret = "CB_API_SECRET"

func main() {

	// Args
	serverPtr := flag.Bool("server", false, "whether or not to start server (default port: 8080)")
	savePtr := flag.Bool("save", true, "whether or not to save a list of transactions")
	transactionTypePtr := flag.String("transaction-type", "all", "the type of coin to parse transactions against")

	// Parse the argument flags
	flag.Parse()

	client := http.Client{
		Timeout: time.Second * 10,
	}
	var absClient query.HTTPClient
	absClient = &client

	fmt.Println("Server enabled:", *serverPtr)
	fmt.Println("Save enabled:", *savePtr)
	fmt.Println("Transaction type:", *transactionTypePtr)

	//filepath := os.Getenv(WarchestConfigEnv)
	//configFile := config.ConfigFile{Filepath: filepath}
	//warchestConfig, err := configFile.ToConfig()
	//if err != nil {
	//	fmt.Printf("Failed loading config: %s\n", err)
	//	os.Exit(FailedLoadConfigRC)
	//}

	apiKey := os.Getenv(CbAPIKey)
	apiSecret := os.Getenv(CbAPISecret)
	cbAuth := auth.CBAuth{apiKey, apiSecret}

	userID, _ := query.CBRetrieveUserID(cbAuth, absClient)
	accountsResp, _ := query.CBRetrieveAccounts(cbAuth, absClient)

	// Filter out non-zero amount of individual coin types
	// Restrict to DOGE for testing reasons. This is the devel branch -P
	valueMap := map[string]query.WarchestCoin{}
	for _, account := range accountsResp.Accounts {
		if account.Balance.Amount > 0.0 &&
			account.Currency.Code == "DOGE" {
			warchestCoin := query.WarchestCoin{AccountID: account.ID, CoinSymbol: account.Currency.Code}

			// Update the coin's rates, profit, and cost
			warchestCoin.Update(cbAuth, absClient)
			valueMap[account.Currency.Code] = warchestCoin
		}
	}

	fmt.Printf("There are %d types of coins.\n", len(accountsResp.Accounts))
	fmt.Printf("You have %d type(s) of coin(s) in your wallet:\n", len(valueMap))
	for _, wcCoin := range valueMap {
		fmt.Printf("\t%s\n", wcCoin.CoinSymbol)
		fmt.Printf("\t\tAmount: %.6f\n", wcCoin.Amount)
		fmt.Printf("\t\tCost: %.6f\n", wcCoin.Cost)
		fmt.Printf("\t\tProfit: %.6f\n", wcCoin.Profit)
	}

	fmt.Printf("\nUpdating crypto wallet for id: %s\n", userID)

	//localWallet := warchestConfig.ToWallet()
	//
	//netProfit, err := wallet.CalculateNetProfit(localWallet, absClient)
	//if err != nil {
	//	fmt.Printf("Failed to calculate Wallet's Profit: %s\n", err)
	//	os.Exit(FailedCalculatingWallet)
	//}
	//
	//fmt.Printf("Current Wallet's Net Profit: %.10f\n", netProfit)

	// Download transactions
	//transactions, err := query.CBCoinTransactions(accountID, cbAuth, absClient)
	//
	//for i, transaction := range transactions {
	//	fmt.Printf("Transaction ID: %s", transaction.Data[i].ID)
	//	fmt.Printf("Transaction ID: %s", transaction.Data[i].Amount.Currency)
	//}
}

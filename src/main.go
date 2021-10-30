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
	// "warchest/src/wallet"
)

const FailedLoadConfigRC = 2
const FailedRetrievingData = 3
const FailedCalculatingWallet = 3
const WarchestConfigEnv = "WARCHEST_CONFIG"
const CbApiKey = "CB_API_KEY"
const CbApiSecret = "CB_API_SECRET"

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
	var absClient query.HttpClient
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

	apiKey := os.Getenv(CbApiKey)
	apiSecret := os.Getenv(CbApiSecret)
	cbAuth := auth.CBAuth{apiKey, apiSecret}

	userID, _ := query.CBRetrieveUserID(cbAuth, absClient)
	accountsResp, _ := query.CBRetrieveAccounts(cbAuth, absClient)

	// Filter out non-zero amount of individual coin types
	valueMap := map[string]query.CBAccount{}
	for _, account := range accountsResp.Accounts {
		if account.Balance.Amount > 0.0 {
			valueMap[account.Currency.Code] = account
		}
	}

	fmt.Printf("There are %d types of coins.\n", len(accountsResp.Accounts))
	fmt.Printf("You have %d type(s) of coin(s) in your wallet.\n", len(valueMap))

	fmt.Printf("Updating crypto wallet for id: %s\n", userID)
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
	//	fmt.Printf("Transaction Id: %s", transaction.Data[i].Id)
	//	fmt.Printf("Transaction Id: %s", transaction.Data[i].Amount.Currency)
	//}
}

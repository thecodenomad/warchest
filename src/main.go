package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
	"warchest/src/auth"
	"warchest/src/config"
	"warchest/src/query"
	"warchest/src/wallet"
)

const FailedLoadConfigRC = 2
const FailedRetrievingData = 3
const FailedCalculatingWallet = 3
const WarchestConfigEnv = "WARCHEST_CONFIG"
const CbApiKey = "CB_API_KEY"
const CbApiSecret = "CB_API_SECRET"

func main() {

	serverPtr := flag.Bool("server", false, "whether or not to start server (default port: 8080")
	client := http.Client{
		Timeout: time.Second * 10,
	}

	fmt.Println("Server enabled:", *serverPtr)

	configPath := os.Getenv(WarchestConfigEnv)

	warchestConfig, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed loading config: %s\n", err)
		os.Exit(FailedLoadConfigRC)
	}

	apiKey := os.Getenv(CbApiKey)
	apiSecret := os.Getenv(CbApiSecret)
	cbAuth := auth.CBAuth{apiKey, apiSecret}

	userId, nil := query.CBRetrieveUserID(cbAuth, client)

	fmt.Printf("Updating crypto wallet for id: %s\n", userId)
	localWallet := warchestConfig.ToWallet()

	netProfit, err := wallet.CalculateNetProfit(localWallet, client)
	if err != nil {
		fmt.Printf("Failed to calculate Wallet's Profit: %s\n", err)
		os.Exit(FailedCalculatingWallet)
	}

	fmt.Printf("Current Wallet's Net Profit: %.10f\n", netProfit)
}

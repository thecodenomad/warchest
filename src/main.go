package main

import (
	"flag"
	"fmt"
	"os"
	"warchest/src/config"
	"warchest/src/wallet"
)

const FailedLoadConfigRC = 2
const FailedRetrievingData = 3
const FailedCalculatingWallet = 3
const WarchestConfigEnv = "WARCHEST_CONFIG"

func main() {

	serverPtr := flag.Bool("server", false, "whether or not to start server (default port: 8080")

	fmt.Println("Server enabled:", *serverPtr)

	configPath := os.Getenv(WarchestConfigEnv)

	warchestConfig, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed loading config: %s\n", err)
		os.Exit(FailedLoadConfigRC)
	}

	fmt.Println("Updating crypto wallet")
	localWallet := warchestConfig.ToWallet()

	netProfit, err := wallet.CalculateNetProfit(localWallet)
	if err != nil {
		fmt.Printf("Failed to calculate Wallet's Profit: %s\n", err)
		os.Exit(FailedCalculatingWallet)
	}

	fmt.Printf("Current Wallet's Net Profit: %.10f\n", netProfit)
}

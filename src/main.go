package main

import (
	"fmt"
	"os"
	"warchest/src/config"
	"warchest/src/query"
	"warchest/src/wallet"
)

const FailedLoadConfigRC = 2
const FailedRetrievingData = 3
const FailedCalculatingWallet = 3
const WarchestConfigEnv = "WARCHEST_CONFIG"

func main() {

	configPath := os.Getenv(WarchestConfigEnv)

	warchestConfig, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed loading config: %s\n", err)
		os.Exit(FailedLoadConfigRC)
	}

	firstCoin := warchestConfig.Transactions[0].CoinSymbol
	fmt.Printf("Warchest config, first coin in config: %s\n", firstCoin)

	coinInfo, err := query.RetrieveCoinData(firstCoin)
	if err != nil {
		fmt.Printf("Failed to retrieve coin data!")
		os.Exit(FailedRetrievingData)
	}
	fmt.Printf("Exchange Rates for %s:\nUSD: %.4f\nGBP: %.4f\nEURO: %.4f\n", firstCoin,
		coinInfo.ExchangeRates.USD,
		coinInfo.ExchangeRates.GBP,
		coinInfo.ExchangeRates.EUR)

	localWallet := warchestConfig.ToWallet()

	netProfit, err := wallet.CalculateNetProfit(localWallet)
	if err != nil {
		fmt.Printf("Failed to calculate Wallet's Profit: %s\n", err)
		os.Exit(FailedCalculatingWallet)
	}

	fmt.Printf("Current Wallet's Net Profit: %.10f\n", netProfit)
}

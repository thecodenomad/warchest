package main

import (
	"fmt"
	"os"
	"warchest/src/config"
	"warchest/src/query"
)

const FailedLoadConfigRC = 2
const FailedRetrievingData = 3
const WarchestConfigEnv = "WARCHEST_CONFIG"

func main() {

	configPath := os.Getenv(WarchestConfigEnv)

	warchestConfig, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed loading config: %s\n", err)
		os.Exit(FailedLoadConfigRC)
	}

	firstCoin := warchestConfig.Coins[0].CoinSymbol
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
}

package main

import (
	"fmt"
	"os"
	"warchest/src/config"
	"warchest/src/query"
)

const FailedLoadConfigRC = 2
const WarchestConfigEnv = "WARCHEST_CONFIG"

func main() {

	configPath := os.Getenv(WarchestConfigEnv)

	warchestConfig, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed loading config: %s\n", err)
		os.Exit(FailedLoadConfigRC)
	}

	firstCoin := warchestConfig.PurchasedCoins[0].CoinSymbol
	fmt.Printf("Warchest config, first coin in config: %s\n", firstCoin)

	coinInfo := query.RetrieveCoinData(firstCoin)
	fmt.Printf("Exchange Rates for %s:\nUSD: %s\nGBP: %s\nEURO: %s\n", firstCoin,
		coinInfo.ExchangeRates.USD,
		coinInfo.ExchangeRates.GBP,
		coinInfo.ExchangeRates.EUR)
}

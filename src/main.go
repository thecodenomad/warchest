package main

import (
	"fmt"
	"os"
	"warchest/src/config"
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

	fmt.Printf("Warchest config, first coin symbol: %s\n", warchestConfig.PurchasedCoins[0].CoinSymbol)
}

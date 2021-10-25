package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"warchest/src/wallet"
)

var (
	ErrReadingFile   = ConfigurationError("Failed reading file!")
	ErrFileNotFound  = ConfigurationError("File not found!")
	ErrMalformedJSON = ConfigurationError("Config isn't correct JSON!")
)

type ConfigurationError string

type Config struct {
	Transactions []Transaction `json:"coin_purchases"`
}

type Transaction struct {
	CoinSymbol        string  `json:"coin_symbol"`
	Amount            float64 `json:"amount"`
	PurchasedPriceUSD float64 `json:"purchased_price_usd"`
	TransactionFee    float64 `json:"transaction_fee"`
}

func (c ConfigurationError) Error() string {
	return string(c)
}

func (c *Config) ToWallet() wallet.Wallet {

	coins := make(map[string]wallet.Coin)

	// Collect coins into a slice
	for _, configTransaction := range c.Transactions {

		coinSymbol := configTransaction.CoinSymbol

		// Is Coin found?
		coin, ok := coins[coinSymbol]
		if !ok {
			coinToInit := wallet.Coin{coinSymbol, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, []wallet.CoinTransaction{}}
			coins[configTransaction.CoinSymbol] = coinToInit

			// Coin to work with for the rest of the transaction collection
			coin = coinToInit
		}

		coinTransaction := wallet.CoinTransaction{configTransaction.Amount, configTransaction.PurchasedPriceUSD, configTransaction.TransactionFee}
		coin.Transactions = append(coin.Transactions, coinTransaction)
		coins[configTransaction.CoinSymbol] = coin
	}

	wallet := wallet.Wallet{[]wallet.Coin{}, 0.0}
	// Convert map to wallet
	for _, coin := range coins {
		// Create new coins from the collection above
		wallet.Coins = append(wallet.Coins, coin)
	}

	return wallet
}

func (t Transaction) ToCoinTransaction() wallet.CoinTransaction {
	return wallet.CoinTransaction{NumCoins: t.Amount, PurchasedPrice: t.PurchasedPriceUSD, TransactionFee: t.TransactionFee}
}

// LoadConfig loads the specified JSON file
func LoadConfig(filepath string) (Config, error) {

	// Check for files existence first
	_, err := os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, ErrFileNotFound
	}

	// Try and load the config
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return Config{}, ErrReadingFile
	}

	// TODO: Properly handle this
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Load the JSON
	tmpConfig := Config{}
	err = json.Unmarshal(byteValue, &tmpConfig)
	if err != nil {
		var synErr *json.SyntaxError
		if errors.As(err, &synErr) {
			return Config{}, ErrMalformedJSON
		}
		return Config{}, ErrReadingFile
	}
	return tmpConfig, err
}

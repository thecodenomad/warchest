package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"warchest/src/query"
	"warchest/src/wallet"
)

var (
	ErrReadingFile   = ConfigurationError("Failed reading file!")
	ErrFileNotFound  = ConfigurationError("File not found!")
	ErrMalformedJSON = ConfigurationError("Config isn't correct JSON!")
	ErrOnUnMarshall  = ConfigurationError("Failed Unmarshalling JSON!")
)

func (c ConfigurationError) Error() string {
	return string(c)
}

type ConfigFile struct {
	Filepath       string
	ByteValue      []byte
	WarchestConfig Config
}

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

func (c *ConfigFile) Exists() bool {
	// Check for files existence first
	_, err := os.Stat(c.Filepath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (c *ConfigFile) Load() ([]byte, error) {
	// Try and load the config
	jsonFile, err := os.Open(c.Filepath)
	if err != nil {
		return []byte{}, ErrReadingFile
	}

	// TODO: Properly handle this
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	c.ByteValue = byteValue

	return byteValue, nil
}

func (c *ConfigFile) Parse() (Config, error) {
	// Unmarshall the JSON
	tmpConfig := Config{}
	err := json.Unmarshal(c.ByteValue, &tmpConfig)
	if err != nil {
		return Config{}, ErrOnUnMarshall
	}
	// Remove c.ByteValue so it's not duplicating space
	c.ByteValue = []byte{}
	return tmpConfig, err
}

func (c *ConfigFile) ToConfig() (Config, error) {

	// Check if the file has already been loaded
	if !c.Exists() {
		return Config{}, ErrFileNotFound
	}

	// Load the byte stream into memory if it hasn't been
	if len(c.ByteValue) == 0 {
		_, err := c.Load()
		if err != nil {
			return Config{}, ErrReadingFile
		}
	}

	// Unmarshall the JSON
	if len(c.WarchestConfig.Transactions) == 0 {
		tmpConfig, err := c.Parse()
		if err != nil {
			return Config{}, ErrOnUnMarshall
		}
		// Don't reload if we don't have to
		c.WarchestConfig = tmpConfig
		return tmpConfig, nil
	}

	return c.WarchestConfig, nil
}

func (c *Config) ToWallet() wallet.Wallet {

	coins := make(map[string]wallet.WarchestCoin)

	// Collect coins into a slice
	for _, configTransaction := range c.Transactions {

		coinSymbol := configTransaction.CoinSymbol

		// Is Coin found?
		coin, ok := coins[coinSymbol]
		if !ok {
			coinToInit := wallet.WarchestCoin{0.0, 0.0, 0.0, query.CoinRates{},
				coinSymbol, []wallet.CoinTransaction{}}

			coins[configTransaction.CoinSymbol] = coinToInit

			// Coin to work with for the rest of the transaction collection
			coin = coinToInit
		}

		coinTransaction := wallet.CoinTransaction{configTransaction.Amount,
			configTransaction.PurchasedPriceUSD, configTransaction.TransactionFee}

		coin.Transactions = append(coin.Transactions, coinTransaction)
		coins[configTransaction.CoinSymbol] = coin
	}

	wallet := wallet.Wallet{[]wallet.WarchestCoin{}, 0.0}
	// Convert map to wallet
	for _, coin := range coins {
		// Create new coins from the collection above
		wallet.Coins = append(wallet.Coins, coin)
	}

	return wallet
}

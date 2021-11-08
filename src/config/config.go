package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"warchest/src/query"
)

var (
	// ErrReadingFile occurs when the config file can't be read
	ErrReadingFile = ConfigurationError("Failed reading file!")

	//ErrFileNotFound occurs when the config file can't be found
	ErrFileNotFound = ConfigurationError("File not found!")

	//ErrMalformedJSON occurs when the JSON for the config file isn't valid
	ErrMalformedJSON = ConfigurationError("Config isn't correct JSON!")

	//ErrOnUnMarshall occurs when the JSON structure doesn't match the configuration structure
	ErrOnUnMarshall = ConfigurationError("Failed Unmarshalling JSON!")
)

// ConfigurationError struct for the errors defined above
type ConfigurationError string

// ConfigurationError helper method to throw the above errors
func (c ConfigurationError) Error() string {
	return string(c)
}

// LocalConfigFile is a wrapper struct around the config file allowing for regeneration and storing of loaded config
type LocalConfigFile struct {
	Filepath       string
	ByteValue      []byte
	WarchestConfig Config
}

// Config is the object that holds transactions pulled in form the config file
type Config struct {
	Transactions []Transaction `json:"coin_purchases"`
}

// Transaction is an individual transaction object used by warchest
type Transaction struct {
	CoinSymbol        string  `json:"coin_symbol"`
	Amount            float64 `json:"amount"`
	PurchasedPriceUSD float64 `json:"purchased_price_usd"`
	TransactionFee    float64 `json:"transaction_fee"`
}

// Exists method that checks if the config file exists
func (c *LocalConfigFile) Exists() bool {
	// Check for files existence first
	_, err := os.Stat(c.Filepath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// Load method loads the config file into a byte array
func (c *LocalConfigFile) Load() ([]byte, error) {
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

// Parse method parses a loaded config file from the internal ByteValue array.
func (c *LocalConfigFile) Parse() (Config, error) {
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

// ToConfig method takes the config file and produces a instantiated config that warchest can consume
func (c *LocalConfigFile) ToConfig() (Config, error) {

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

// ToWallet method that produces a wallet based on the config object
func (c *Config) ToWallet() query.Wallet {

	coins := make(map[string]query.WarchestCoin)

	// Collect coins into a slice
	for _, configTransaction := range c.Transactions {

		coinSymbol := configTransaction.CoinSymbol

		// Is Coin found?
		coin, ok := coins[coinSymbol]
		if !ok {
			coinToInit := query.WarchestCoin{"", 0.0, 0.0, 0.0, query.CoinRates{},
				coinSymbol, []query.CoinTransaction{}}

			coins[configTransaction.CoinSymbol] = coinToInit

			// Coin to work with for the rest of the transaction collection
			coin = coinToInit
		}

		coinTransaction := query.CoinTransaction{configTransaction.Amount,
			configTransaction.PurchasedPriceUSD, configTransaction.TransactionFee}

		coin.Transactions = append(coin.Transactions, coinTransaction)
		coins[configTransaction.CoinSymbol] = coin
	}

	wallet := query.Wallet{map[string]query.WarchestCoin{}, 0.0}
	// Convert map to wallet
	for _, coin := range coins {
		// Create new coins from the collection above
		wallet.Coins[coin.Symbol] = coin
	}

	return wallet
}

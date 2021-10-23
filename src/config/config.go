package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type ConfigErr string

var (
	ErrReadingFile   = ConfigErr("Failed reading file!")
	ErrFileNotFound  = ConfigErr("File not found!")
	ErrMalformedJSON = ConfigErr("Config isn't correct JSON!")
)

func (c ConfigErr) Error() string {
	return string(c)
}

type Config struct {
	PurchasedCoins []struct {
		CoinSymbol           string  `json:"coin_symbol"`
		Amount               float64 `json:"amount"`
		PurchasedPrice       float64 `json:"purchased_price"`
		TransactionFee       float64 `json:"transaction_fee"`
		PurchaseExchangeRate float64 `json:"purchase_exchange_rate"`
	} `json:"purchased_coins"`
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

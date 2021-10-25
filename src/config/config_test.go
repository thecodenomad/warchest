package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Test Valid Config", func(t *testing.T) {
		tmpConfig, err := LoadConfig("./testdata/CoinConfig.json")
		assert.Nil(t, err, "Should not fail loading string")
		assert.Equal(t, len(tmpConfig.Transactions), 2)
		assert.Equal(t, tmpConfig.Transactions[0].CoinSymbol, "ETH")

		valueTests := []struct {
			actualValue   float64
			expectedValue float64
		}{
			{tmpConfig.Transactions[0].Amount, 10.1},
			{tmpConfig.Transactions[0].PurchasedPriceUSD, 34.5},
			{tmpConfig.Transactions[0].TransactionFee, 6.56},
			{tmpConfig.Transactions[1].Amount, 5.0},
			{tmpConfig.Transactions[1].PurchasedPriceUSD, 2.5},
			{tmpConfig.Transactions[1].TransactionFee, 0.35},
		}

		// Validate the rest of the imported values
		for _, tt := range valueTests {
			assert.Equal(t, tt.expectedValue, tt.actualValue)
		}
	})

	t.Run("Test non-existent filepath", func(t *testing.T) {
		tmpConfig, err := LoadConfig("./testdata/Bogus.json")
		assert.Equal(t, tmpConfig, Config{}, "Config not empty!")
		assert.Error(t, err, "should have raised an error")
		assert.Equal(t, ErrFileNotFound, err, "should have raised a file not found error")
	})

	//t.Run("Test permissions issues (simulate corruption)", func(t *testing.T) {
	//	tmpConfig, err := LoadConfig("./testdata/Unreadable.json")
	//	assert.Equal(t, tmpConfig, Config{}, "Config not empty!")
	//	assert.Error(t, err, "should have raised an error")
	//	assert.Equal(t, ErrReadingFile, err, "should have raised a read error")
	//})

	t.Run("Test malformed JSON", func(t *testing.T) {
		tmpConfig, err := LoadConfig("./testdata/Malformed.json")
		assert.Equal(t, tmpConfig, Config{}, "Config not empty!")
		assert.Error(t, err, "should have raised an error")
		assert.Equal(t, ErrOnUnMarshall, err, "should have raised a malformed error")
	})

	t.Run("Convert config to wallet", func(t *testing.T) {
		tmpConfig, _ := LoadConfig("./testdata/CoinConfig.json")
		wallet := tmpConfig.ToWallet()
		assert.Equal(t, len(wallet.Coins), 2, "Failed to have correct number of coins")

		// Each coin should have 1 transaction
		for _, coin := range wallet.Coins {
			assert.Equal(t, len(coin.Transactions), 1, "Failed to have correct number of transactions")
		}
	})

	// addmittingly overkill, and borderline useful, but it makes for full coverage!
	t.Run("Test error handling", func(t *testing.T) {
		valueTests := []struct {
			actualValue   string
			expectedValue string
		}{
			{ErrReadingFile.Error(), "Failed reading file!"},
			{ErrFileNotFound.Error(), "File not found!"},
			{ErrMalformedJSON.Error(), "Config isn't correct JSON!"},
			{ErrOnUnMarshall.Error(), "Failed Unmarshalling JSON!"},
		}
		// Validate the rest of the imported values
		for _, tt := range valueTests {
			assert.Equal(t, tt.expectedValue, tt.actualValue)
		}
	})
}

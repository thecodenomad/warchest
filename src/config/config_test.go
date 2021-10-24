package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Test Valid Config", func(t *testing.T) {
		tmpConfig, err := LoadConfig("./testdata/CoinConfig.json")
		assert.Nil(t, err, "Should not fail loading string")
		assert.Equal(t, len(tmpConfig.Coins), 2)
		assert.Equal(t, tmpConfig.Coins[0].CoinSymbol, "ETH")

		valueTests := []struct {
			actualValue   float64
			expectedValue float64
		}{
			{tmpConfig.Coins[0].Amount, 10.1},
			{tmpConfig.Coins[0].PurchasedPrice, 34.5},
			{tmpConfig.Coins[0].TransactionFee, 6.56},
			{tmpConfig.Coins[0].PurchaseRateUSD, 0.001},
			{tmpConfig.Coins[1].Amount, 5.0},
			{tmpConfig.Coins[1].PurchasedPrice, 2.5},
			{tmpConfig.Coins[1].TransactionFee, 0.35},
			{tmpConfig.Coins[1].PurchaseRateUSD, 0.40},
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

	t.Run("Test malformed JSON", func(t *testing.T) {
		tmpConfig, err := LoadConfig("./testdata/Malformed.json")
		assert.Equal(t, tmpConfig, Config{}, "Config not empty!")
		assert.Error(t, err, "should have raised an error")
		assert.Equal(t, ErrMalformedJSON, err, "should have raised a malformed error")
	})
}

package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {

	// Happy Path
	t.Run("Test Valid Config", func(t *testing.T) {

		testConfigFile := LocalFile{Filepath: "./testdata/CoinConfig.json"}
		tmpConfig, err := testConfigFile.ToConfig()

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

		newConfig, err := testConfigFile.ToConfig()
		assert.NotNil(t, newConfig, "should have been same config as initially loaded")
	})

	t.Run("Convert config to wallet", func(t *testing.T) {
		testConfigFile := LocalFile{Filepath: "./testdata/CoinConfig.json"}
		tmpConfig, _ := testConfigFile.ToConfig()
		wallet := tmpConfig.ToWallet()
		assert.Equal(t, len(wallet.Coins), 2, "Failed to have correct number of coins")

		// Each coin should have 1 transaction
		for _, coin := range wallet.Coins {
			assert.Equal(t, len(coin.Transactions), 1, "Failed to have correct number of transactions")
		}
	})

	// Cloudy Path
	t.Run("Test file existence", func(t *testing.T) {

		// Test non existent file
		testConfigFile := LocalFile{Filepath: "./testdata/Bogus.json"}
		assert.False(t, testConfigFile.Exists(), "This file should not exist")

		// Test existing file
		testConfigFile = LocalFile{Filepath: "./testdata/Malformed.json"}
		assert.True(t, testConfigFile.Exists(), "This file should exist")

	})

	t.Run("Test malformed JSON", func(t *testing.T) {
		testConfigFile := LocalFile{Filepath: "./testdata/Malformed.json"}
		tmpConfig, err := testConfigFile.Parse()

		assert.Equal(t, tmpConfig, Config{}, "Config not empty!")
		assert.Error(t, err, "should have raised an error")
		assert.Equal(t, ErrOnUnMarshall, err, "should have raised a malformed error")
	})

	t.Run("Test file read problems", func(t *testing.T) {
		testConfigFile := LocalFile{Filepath: "./testdata/Bogus.json"}
		byteValue, err := testConfigFile.Load()

		assert.Empty(t, byteValue, "should be empty")
		assert.Error(t, err, "should have raised an error")
		assert.Equal(t, ErrReadingFile, err, "should have raised a read error")
	})

	t.Run("Test To Config error handling", func(t *testing.T) {

		// Existence check
		testConfigFile := LocalFile{Filepath: "./testdata/Bogus.json"}
		_, err := testConfigFile.ToConfig()
		assert.NotNil(t, ErrFileNotFound, err, "this should have been a file existence error")

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

//
// Mocks
///////////////////////

type MockConfigParser interface {
	Load() ([]byte, error)
	Parse() (Config, error)
}

type MockConfigLoader interface {
	ToConfig() (Config, error)
}

type MockConfigReader interface {
	Load() ([]byte, error)
}

type MockConfigFile struct {
	Filepath string
}

func (m *MockConfigFile) Exits() bool {
	return false
}

func (m *MockConfigFile) Load() ([]byte, error) {
	return nil, ErrReadingFile
}

func (m *MockConfigFile) ToConfig() (Config, error) {
	_, err := m.Load()
	if err != nil {
		return Config{}, ErrReadingFile
	}
	return Config{}, ErrReadingFile
}

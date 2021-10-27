package query

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
	"warchest/src/auth"
)

// TODO: should be using a mock (spy) so we aren't making http requests
func TestRetrieveCoinData(t *testing.T) {
	// Test variables
	symbol := "ETH"

	t.Run("Happy Path", func(t *testing.T) {
		coin := "BTC"
		coinInfo, err := CBRetrieveCoinData(coin)
		assert.Nil(t, err, "failed to retrieve rates")
		assert.Equal(t, coin, coinInfo.Currency, "values should be the same!")
		assert.NotEqual(t, CoinInfo{}, coinInfo, "coinInfo object should not be empty")
		assert.NotNil(t, coinInfo.ExchangeRates, "no rates found!")
	})

	t.Run("Rainy Day connectivity!", func(t *testing.T) {

		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBExchangeRateUrl).
			Reply(500).
			BodyString(`nada`)

		coinInfo, err := CBRetrieveCoinData(symbol)

		// There should have been a connection error
		assert.NotNil(t, err, "This call should have produced a connection error")
		assert.Equal(t, CoinInfo{}, coinInfo)
	})

	t.Run("Malformed JSON response", func(t *testing.T) {

		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBExchangeRateUrl).
			Reply(200).
			BodyString(`[asdf,[],!}`)

		coinInfo, err := CBRetrieveCoinData(symbol)

		// There should have been a connection error
		assert.NotNil(t, err, "This call should have produced a JSON parse error")
		assert.Equal(t, CoinInfo{}, coinInfo)
	})

	// TODO: Still having problems forcing failure with io.ReadAll need to look into this
	t.Run("Failure to read response body check", func(t *testing.T) {

		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBExchangeRateUrl).
			Reply(200).
			SetHeader("Content-Length", "10")

		coinInfo, err := CBRetrieveCoinData(symbol)

		assert.NotNil(t, err, "This call should have produced a read error for the response body")
		assert.Equal(t, CoinInfo{}, coinInfo)
	})

}

func TestCBRetrieveUserID(t *testing.T) {

	// Setup Test specifics
	cbAuth := auth.CBAuth{}
	json := `{ "data": { "id": "9da7a204-544e-5fd1-9a12-61176c5d4cd8" } }`

	// Estbalish Mock
	defer gock.Off()
	gock.New(CBBaseURL).
		Get(CBUserUrl).
		Reply(200).
		BodyString(json)

	actualResp, _ := CBRetrieveUserID(cbAuth)
	expectedResp := "9da7a204-544e-5fd1-9a12-61176c5d4cd8"
	assert.Equal(t, expectedResp, actualResp)
}

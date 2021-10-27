package query

import (
	"github.com/jarcoal/httpmock"
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

		json := `{"data": { "currency": "ETH", "rates": {"USD": "12.0", "EUR": "11.0", "GBP": "10.0"}}}`
		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBExchangeRateUrl).
			Reply(200).
			BodyString(json)

		coinInfo, err := CBRetrieveCoinData(symbol)
		assert.Nil(t, err, "failed to retrieve rates")
		assert.Equal(t, symbol, coinInfo.Currency, "values should be the same!")
		assert.Equal(t, 12.00, coinInfo.ExchangeRates.USD, "Should be the same")
		assert.Equal(t, 11.00, coinInfo.ExchangeRates.EUR, "Should be the same")
		assert.Equal(t, 10.00, coinInfo.ExchangeRates.GBP, "Should be the same")
		assert.NotEqual(t, CoinInfo{}, coinInfo, "coinInfo object should not be empty")
		assert.NotNil(t, coinInfo.ExchangeRates, "no rates found!")
	})

	t.Run("Rainy Day connectivity!", func(t *testing.T) {

		// Easier to do with httpmock, default behavior is to return a connection error
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

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

	t.Run("Happy Path", func(t *testing.T) {
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
	})

	t.Run("Cloudy Path", func(t *testing.T) {
		// Setup Test specifics
		cbAuth := auth.CBAuth{}
		json := `{unparseable,,,}`

		// Estbalish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBUserUrl).
			Reply(200).
			BodyString(json)

		actualResp, _ := CBRetrieveUserID(cbAuth)
		expectedResp := ""
		assert.Equal(t, expectedResp, actualResp)
	})
}

package query

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: should be using a mock (spy) so we aren't making http requests
func TestRetrieveCoinData(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		coin := "BTC"
		coinInfo, err := CBRetrieveCoinData(coin)
		assert.Nil(t, err, "failed to retrieve rates")
		assert.Equal(t, coin, coinInfo.Currency, "values should be the same!")
		assert.NotEqual(t, CoinInfo{}, coinInfo, "coinInfo object should not be empty")
		assert.NotNil(t, coinInfo.ExchangeRates, "no rates found!")
	})

	t.Run("Rainy Day connectivity!", func(t *testing.T) {

		// Test variables
		symbol := "ETH"

		// Setup capture for HTTP call
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		coinInfo, err := CBRetrieveCoinData(symbol)

		// Stop mock
		httpmock.GetTotalCallCount()

		// There should have been a connection error
		assert.NotNil(t, err, "This call should have produced a connection error")

		// Empty CoinInfo object at this point
		assert.Equal(t, CoinInfo{}, coinInfo)
	})

	t.Run("Malformed JSON response", func(t *testing.T) {

		// Test variables
		symbol := "ETH"
		url := CBExchangeRateUrl + "?currency=" + symbol

		// Setup capture for HTTP call
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		mockedResponse := httpmock.NewStringResponder(200,
			`{bad json and stuff}`)
		httpmock.RegisterResponder("GET", url, mockedResponse)

		coinInfo, err := CBRetrieveCoinData(symbol)

		// Stop mock
		httpmock.GetTotalCallCount()

		// There should have been a connection error
		assert.NotNil(t, err, "This call should have produced a JSON parse error")

		// Empty CoinInfo object at this point
		assert.Equal(t, CoinInfo{}, coinInfo)
	})
}

//func TestCBRetrieveUserID(t *testing.T) {
//
//	cbAuth := auth.CBAuth{}
//	url := CBUserUrl + CBUserUrl
//
//	httpmock.Activate()
//	defer httpmock.DeactivateAndReset()
//	mockedResponse := httpmock.NewStringResponder(200,
//		`	{ "data": { "id": "9da7a204-544e-5fd1-9a12-61176c5d4cd8" } }`)
//	httpmock.RegisterResponder("GET", url, mockedResponse)
//
//	actualResp, _ := CBRetrieveUserID(cbAuth)
//	expectedResp := "9da7a204-544e-5fd1-9a12-61176c5d4cd8"
//	assert.Equal(t, expectedResp, actualResp)
//}

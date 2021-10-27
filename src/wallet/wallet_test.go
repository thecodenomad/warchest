package wallet

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"warchest/src/query"
)

func TestCalculateCoinProfit(t *testing.T) {

	transactions := []CoinTransaction{
		{1.0, 2.0, 3.0}, // $5 cost
		{4.0, 5.0, 6.0}, // $26 cost
	}

	testRateUSD := 30.0
	expectedNetProfit := 5*testRateUSD - (1.0*2.0 + 3.0 + 4.0*5.0 + 6.0)
	testCoin := Coin{"ETH", testRateUSD, 40.0, 50.0,
		25, 5.0, 31.0, 0.0, transactions}

	testCoin.UpdateProfit()

	assert.Equal(t, expectedNetProfit, testCoin.Profit, "should be the same")
}

// This is one function purely for the coverage stats ;) for 'Update' method
func TestCoin_Update(t *testing.T) {

	// Test variables
	symbol := "ETH"
	url := query.CBBaseURL + query.CBExchangeRateUrl + "?currency=" + symbol

	// Setup capture for HTTP call
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Establish Mocked Response
	mockedResponse := httpmock.NewStringResponder(200,
		`{"data":{"currency":"ETH","rates":{"USD":"12.99","EUR":"11.99","GBP": "10.99"}}}`)
	httpmock.RegisterResponder("GET", url, mockedResponse)

	testAmount := 1.0
	testCost := 10.0
	testFee := 1.0
	testTransactions := []CoinTransaction{{testAmount, testCost, testFee}}
	coin := Coin{"ETH", 0.0, 0.0, 0.0, 0.0,
		0, 0.0, 0.0, testTransactions}

	// Set Expectations for dem noty bits
	expectedRate := 12.99
	expectedCost := testCost*testAmount + testFee
	expectedProfit := expectedRate*testAmount - expectedCost

	// Do the thing (ie. run the 3 update methods)
	coin.Update()

	// Stop mock
	httpmock.GetTotalCallCount()

	// Make sure the algo translated the response correctly
	assert.Equal(t, expectedRate, coin.CurrentRateUSD, "should be the same")
	assert.Equal(t, expectedCost, coin.Cost, "should be the same")
	assert.Equal(t, expectedProfit, coin.Profit, "should be the same")

	// Make sure there was only 1 call to the remote API, we don't want to be banned!
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["GET "+url], "should be the same")

}

func TestCoin_UpdateRates_Cloudy(t *testing.T) {

	// Setup capture for HTTP call
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Establish Mocked Response
	// NOTE: no mocked responses, without registering a response for a url httpmock returns an error by default
	expectedResp := 0.0
	testTransactions := []CoinTransaction{}
	coin := Coin{"ETH", 12.0, 11.0, 10.0, 0.0,
		0, 0.0, 0.0, testTransactions}

	// Update the rates, but since there is an error we should _silently_ ignore and leave the rate at 0
	// TODO: better error handling around requests maybe needed
	coin.UpdateRates()

	// Stop mock
	httpmock.GetTotalCallCount()

	// Verify method corralled the bits
	assert.Equal(t, expectedResp, coin.CurrentRateUSD, "should be the same")
}

func TestCalculateNetProfit(t *testing.T) {

	// Test variables
	symbol := "ETH"
	url := query.CBBaseURL + query.CBExchangeRateUrl + "?currency=" + symbol

	// Pedantic to help with extensibility
	testAmount := 1.0
	testCost := 10.0
	testFee := 1.0
	testTransactions := []CoinTransaction{{testAmount, testCost, testFee}}
	coin := Coin{"ETH", 0.0, 0.0, 0.0, 0.0,
		0, 0.0, 0.0, testTransactions}
	wallet := Wallet{[]Coin{coin}, 0.0}

	// Criteria
	expectedProfit := "1.990000"

	// Setup capture for HTTP call
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// Establish Response
	mockedResponse := httpmock.NewStringResponder(200,
		`{"data":{"currency":"ETH","rates":{"USD":"12.99","EUR":"11.99","GBP": "10.99"}}}`)
	httpmock.RegisterResponder("GET", url, mockedResponse)

	// Do the things then set threshold for easier comparison of float values
	actualResp, err := CalculateNetProfit(wallet)
	actualProfit := fmt.Sprintf("%.6f", actualResp)

	// Stop http capture
	httpmock.GetTotalCallCount()

	// Make sure there was only 1 call to the remote API, we don't want to be banned!
	info := httpmock.GetCallCountInfo()
	assert.Nil(t, err, "this was mocked, and should not fail")
	assert.Equal(t, expectedProfit, actualProfit, "should be the same")
	assert.Equal(t, 1, info["GET "+url], "There should have only been 1 call for 1 coin's rate")
}

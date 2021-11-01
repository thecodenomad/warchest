package query

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
	"time"
	"warchest/src/auth"
)

func TestCoinUpdateProfit(t *testing.T) {

	symbol := "ETH"
	testRateUSD := 30.0
	testCoin := WarchestCoin{"somethingLong", 50.0, 5.0,
		0.0, CoinRates{0.0, 0.0, testRateUSD}, symbol, []CoinTransaction{}}
	expectedNetProfit := 5*testRateUSD - 50

	testCoin.UpdateProfit()

	assert.Equal(t, expectedNetProfit, testCoin.Profit, "should be the same")
}

// This is one function purely for the coverage stats ;) for 'Update' method
func TestCoin_Update(t *testing.T) {

	symbol := "ETH"
	client := http.Client{
		Timeout: time.Second * 10,
	}

	var absClient HTTPClient
	absClient = &client

	// Establish Mock
	json := `{"data":{"currency":"ETH","rates":{"USD":"12.99","EUR":"11.99","GBP": "10.99"}}}`
	defer gock.Off()
	gock.New(CBBaseURL).
		Get(CBExchangeRateURL).
		Reply(200).
		BodyString(json)

	testAmount := 1.0
	testCost := 10.0
	testFee := 1.0
	testRateUSD := 30.0
	accountID := "somethingLong"
	testTransactions := []CoinTransaction{{testAmount, testCost, testFee}}
	testCoin := WarchestCoin{"somethingLong", 5.0, 0.0,
		0.0, CoinRates{0.0, 0.0, testRateUSD}, symbol, testTransactions}

	transactionURL := "/v2/accounts/" + accountID + "/transactions"
	fmt.Printf("Transaction URL to mock: %s\n", transactionURL)
	gock.New(CBBaseURL).
		Get(transactionURL).
		Reply(200).
		BodyString(transactionJSON)

	// Set Expectations for dem noty bits
	expectedRate := 12.99
	expectedCost := testCost*testAmount + testFee
	expectedProfit := expectedRate*testAmount - expectedCost

	// Do the thing (ie. run the 3 update methods)
	testCoin.Update(auth.CBAuth{}, absClient)

	// Make sure the algo translated the response correctly
	assert.Equal(t, expectedRate, testCoin.Rates.USD, "should be the same")
	assert.Equal(t, expectedCost, testCoin.Cost, "should be the same")
	assert.Equal(t, expectedProfit, testCoin.Profit, "should be the same")
}

func TestCoin_UpdateRates_Cloudy(t *testing.T) {

	symbol := "ETH"
	mockClient := &MockClient{}

	// Force request failure
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expectedResp := 0.0
	testTransactions := []CoinTransaction{}
	testCoin := WarchestCoin{"somethingLong", 5.0, 0.0,
		0.0, CoinRates{USD: -10.0}, symbol, testTransactions}

	// Update the rates, but since there is an error we should _silently_ ignore and leave the rate at 0
	// TODO: better error handling around requests maybe needed
	testCoin.UpdateRates(mockClient)

	// Verify method corralled the bits
	assert.Equal(t, expectedResp, testCoin.Rates.USD, "should be the same")
}

func TestCalculateNetProfit(t *testing.T) {

	symbol := "ETH"
	client := http.Client{
		Timeout: time.Second * 10,
	}
	var absClient HTTPClient
	absClient = &client

	// Test variables, pedantic for extensibility
	testAmount := 1.0
	testCost := 10.0
	testFee := 1.0
	accountID := "somethingLong"
	testTransactions := []CoinTransaction{{testAmount, testCost, testFee}}
	testCoin := WarchestCoin{"somethingLong", 5.0, 0.0,
		0.0, CoinRates{USD: -10.0}, symbol, testTransactions}

	wallet := Wallet{[]WarchestCoin{testCoin}, 0.0}

	// Criteria
	expectedProfit := "1.990000"

	// Establish Mock for Exchange Rate
	exchangeJSON := `{"data":{"currency":"ETH","rates":{"USD":"12.99","EUR":"11.99","GBP": "10.99"}}}`
	defer gock.Off()
	gock.New(CBBaseURL).
		Get(CBExchangeRateURL).
		Reply(200).
		BodyString(exchangeJSON)

	// Establish Mock for updating transactions
	transactionURL := "/v2/accounts/" + accountID + "/transactions"
	fmt.Printf("Transaction URL to mock: %s\n", transactionURL)
	gock.New(CBBaseURL).
		Get(transactionURL).
		Reply(200).
		BodyString(transactionJSON)

	// Do the things then set threshold for easier comparison of float values
	actualResp, err := CalculateNetProfit(wallet, auth.CBAuth{}, absClient)
	actualProfit := fmt.Sprintf("%.6f", actualResp)

	// Make sure there was only 1 call to the remote API, we don't want to be banned!
	assert.Nil(t, err, "this was mocked, and should not fail")
	assert.Equal(t, expectedProfit, actualProfit, "should be the same")
}

const transactionJSON = `{
	"pagination": {
		"ending_before": null,
		"starting_after": null,
		"limit": 25,
		"order": "desc",
		"previous_uri": null,
		"next_uri": null
	},
	"data": [
		{
			"id": "4117f7d6-5694-5b36-bc8f-847509850ea4",
			"type": "buy",
			"status": "completed",
			"amount": {
				"amount": "1.00",
				"currency": "ETH"
			},
			"native_amount": {
				"amount": "11.00",
				"currency": "USD"
			},
			"description": null,
			"created_at": "2015-03-26T23:44:08-07:00",
			"updated_at": "2015-03-26T23:44:08-07:00",
			"resource": "transaction",
			"resource_path": "/v2/accounts/somethingLong/transactions/4117f7d6-5694-5b36-bc8f-847509850ea4",
			"buy": {
				"id": "9e14d574-30fa-5d85-b02c-6be0d851d61d",
				"resource": "buy",
				"resource_path": "/v2/accounts/2bbf394c-193b-5b2a-9155-3b4732659ede/buys/9e14d574-30fa-5d85-b02c-6be0d851d61d"
			},
			"details": {
				"title": "Bought ETH",
				"subtitle": "Used some account"
			}
		}
	]
}
`

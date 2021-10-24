package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: should be using a mock (spy) so we aren't making http requests
func TestRetrieveCoinData(t *testing.T) {
	coin := "BTC"
	coinInfo, err := RetrieveCoinData(coin)
	assert.Nil(t, err, "failed to retrieve rates")
	assert.Equal(t, coin, coinInfo.Currency, "values should be the same!")
	assert.NotEqual(t, CoinInfo{}, coinInfo, "coinInfo object should not be empty")
	assert.NotNil(t, coinInfo.ExchangeRates, "no rates found!")
}

func TestCalculateCoinProfit(t *testing.T) {

	transactions := []CoinTransaction{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
	}

	testCoin := Coin{"ETH", 30.0, 40.0, 50.0, transactions}
	coinInfo := CoinInfo{"ETH", Rates{7.0, 6.0, 25.0}}

	// Current USD Rate * Number of Coins - Initial expense
	expectedResponse := coinInfo.ExchangeRates.USD*5.0 - 31.0
	actualResp, err := CalculateCoinProfit(testCoin, coinInfo)

	assert.Nil(t, err, "failed to calculate profits")
	assert.Equal(t, expectedResponse, actualResp, "should be the same")

}

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

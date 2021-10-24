package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: should be using a mock (spy) so we aren't making http requests
func TestRetrieveCoinData(t *testing.T) {
	coin := "BTC"
	coinInfo := RetrieveCoinData(coin)
	assert.Equal(t, coin, coinInfo.Currency, "values should be the same!")
	assert.NotEqual(t, CoinInfo{}, coinInfo, "coreInfo object should not be empty")
	assert.NotNil(t, coinInfo.ExchangeRates, "no rates found!")
}

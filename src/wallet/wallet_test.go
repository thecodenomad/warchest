package wallet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateCoinProfit(t *testing.T) {

	transactions := []CoinTransaction{
		{1.0, 2.0, 3.0}, // $5 cost
		{4.0, 5.0, 6.0}, // $26 cost
	}

	testRateUSD := 30.0
	expectedNetProfit := 5*testRateUSD - (1.0*2.0 + 3.0 + 4.0*5.0 + 6.0)
	testCoin := Coin{"ETH", testRateUSD, 40.0, 50.0, 25, 5.0, 31.0, 0.0, transactions}

	testCoin.UpdateProfit()

	assert.Equal(t, expectedNetProfit, testCoin.Profit, "should be the same")
}

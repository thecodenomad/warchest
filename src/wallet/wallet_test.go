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

	testCoin := Coin{"ETH", 30.0, 40.0, 50.0, 25, transactions}
	expectedNetProfit := 5*testCoin.CurrentRateUSD - (1.0*2.0 + 3.0 + 4.0*5.0 + 6.0)

	actualResp, err := CalculateCoinProfit(testCoin)

	assert.Nil(t, err, "failed to calculate profits")
	assert.Equal(t, expectedNetProfit, actualResp, "should be the same")

}

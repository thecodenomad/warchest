package config

type CryptoSymbol string

type PurchasedCoin struct {
	coin                   CryptoSymbol
	amount                 float64
	purchasedPrice         float64
	transactionFee         float64
	exchangeRateAtPurchase float64
}

type Config struct {
	PurchasedCoins []struct {
		CoinSymbol           string  `json:"coin_symbol"`
		Amount               float64 `json:"amount"`
		PurchasedPrice       float64 `json:"purchased_price"`
		TransactionFee       float64 `json:"transaction_fee"`
		PurchaseExchangeRate float64 `json:"purchase_exchange_rate"`
	} `json:"purchased_coins"`
}

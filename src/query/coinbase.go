package query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"warchest/src/auth"
)

const CBBaseURL = "https://api.coinbase.com"
const CBExchangeRateUrl = "/v2/exchange-rates"
const CBUserUrl = "/v2/user"
const TransactionUrl = "/v2/accounts/:account_id/transactions"

type CoinInfoResponse struct {
	Data CoinInfo `json:"data"`
}

type CoinInfo struct {
	Currency      string `json:"currency"`
	ExchangeRates Rates  `json:"rates"`
}

type Rates struct {
	EUR float64 `json:"EUR,string"`
	GBP float64 `json:"GBP,string"`
	USD float64 `json:"USD,string"`
}

type CBTransaction struct {
}

type CBUserResponse struct {
	Data struct {
		Id              string      `json:"id"`
		Name            string      `json:"name"`
		Username        string      `json:"username"`
		ProfileLocation interface{} `json:"profile_location"`
		ProfileBio      interface{} `json:"profile_bio"`
		ProfileUrl      string      `json:"profile_url"`
		AvatarUrl       string      `json:"avatar_url"`
		Resource        string      `json:"resource"`
		ResourcePath    string      `json:"resource_path"`
	} `json:"data"`
}

// CBRetrieveCoinData will return exchange rates for a given Crypto Currency Symbol
// TODO: Change the structure so that this is a struct method
func CBRetrieveCoinData(symbol string) (CoinInfo, error) {
	url := CBBaseURL + CBExchangeRateUrl + "?currency=" + symbol
	resp, err := http.Get(url)

	//TODO: Create custom error for connectivity failures
	if err != nil {
		log.Printf("Hit error on retrieval: %s", err)
		return CoinInfo{}, err
	}
	defer resp.Body.Close()

	cResp := CoinInfoResponse{}

	//TODO: Create custom error for failures to read respBody as a string
	bodyAsStr, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body of error: %s", err)
		return CoinInfo{}, err
	}

	//TODO: Create custom error for unmarshalling issues
	if err := json.Unmarshal([]byte(bodyAsStr), &cResp); err != nil {
		log.Printf("Failed to unmarshall the bits: %s", err)
		return CoinInfo{}, err
	}

	return cResp.Data, err
}

func CBRetrieveUserID(cbAuth auth.CBAuth) (string, error) {

	url := CBBaseURL + CBUserUrl
	log.Printf("URL: " + url)

	authHeaders := cbAuth.NewAuthMap("GET", "", CBUserUrl)
	req, err := http.NewRequest("GET", url, nil)

	// Set auth headers
	for key, value := range authHeaders {
		log.Printf("Setting %s to %s", key, value)
		req.Header.Add(key, value)
	}

	// Retrieve response
	resp, err := http.DefaultClient.Do(req)

	//TODO: Create custom error for failure to decode
	if err != nil {
		log.Printf("URL: " + url)
		log.Printf("%s", err)
		return "", err
	}
	defer resp.Body.Close()

	userResp := CBUserResponse{}

	//TODO: Create custom error for failure to decode
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		log.Printf("error: %s", err)
		return "", err
	}

	return userResp.Data.Id, nil
}

// CBRetrieveTransactions will return transactions for all coins the apikey has access to
//func CBRetrieveTransactions(auth auth.CBAuth) []CBTransaction {
//
//	return []CBTransaction{}
//}

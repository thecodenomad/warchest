package query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"warchest/src/auth"
)

const CBBaseURL = "https://api.coinbase.com"
const CBExchangeRateUrl = "/v2/exchange-rates"
const CBUserUrl = "/v2/user"
const CBTransactionUrl = "/v2/accounts/:account_id/transactions"

var (
	ErrDecoding     = QueryError("failed decoding response")
	ErrOnUnmarshall = QueryError("failed to unmarshall")
	ErrConnection   = QueryError("error during request")
)

func (q QueryError) Error() string {
	return string(q)
}

type QueryError string

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
	Pagination struct {
		EndingBefore  interface{} `json:"ending_before"`
		StartingAfter interface{} `json:"starting_after"`
		Limit         int         `json:"limit"`
		Order         string      `json:"order"`
		PreviousUri   interface{} `json:"previous_uri"`
		NextUri       interface{} `json:"next_uri"`
	} `json:"pagination"`
	Data []struct {
		Id     string `json:"id"`
		Type   string `json:"type"`
		Status string `json:"status"`
		Amount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"amount"`
		NativeAmount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"native_amount"`
		Description  *string   `json:"description"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Resource     string    `json:"resource"`
		ResourcePath string    `json:"resource_path"`
		Buy          struct {
			Id           string `json:"id"`
			Resource     string `json:"resource"`
			ResourcePath string `json:"resource_path"`
		} `json:"buy,omitempty"`
		Details struct {
			Title    string `json:"title"`
			Subtitle string `json:"subtitle"`
		} `json:"details"`
		To struct {
			Resource     string `json:"resource"`
			Email        string `json:"email,omitempty"`
			Id           string `json:"id,omitempty"`
			ResourcePath string `json:"resource_path,omitempty"`
		} `json:"to,omitempty"`
		Network struct {
			Status string `json:"status"`
			Name   string `json:"name"`
		} `json:"network,omitempty"`
	} `json:"data"`
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

// HttpClient interface is an internal interface useful for testability
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CBRetrieveCoinData will return exchange rates for a given Crypto Currency Symbol
// TODO: Change the structure so that this is a struct method
func CBRetrieveCoinData(symbol string, client HttpClient) (CoinInfo, error) {
	url := CBBaseURL + CBExchangeRateUrl + "?currency=" + symbol

	req, err := http.NewRequest("GET", url, nil)

	// Retrieve response
	resp, err := client.Do(req)

	//TODO: Create custom error for connectivity failures
	if err != nil {
		log.Printf("Hit error on retrieval: %s", err)
		return CoinInfo{}, ErrConnection
	}
	defer resp.Body.Close()

	cResp := CoinInfoResponse{}

	//TODO: Create custom error for failures to read respBody as a string
	bodyAsStr, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body of error: %s", err)
		return CoinInfo{}, ErrDecoding
	}

	//TODO: Create custom error for unmarshalling issues
	if err := json.Unmarshal([]byte(bodyAsStr), &cResp); err != nil {
		log.Printf("%s", err)
		return CoinInfo{}, ErrOnUnmarshall
	}

	return cResp.Data, err
}

func CBRetrieveUserID(cbAuth auth.CBAuth, client HttpClient) (string, error) {

	url := CBBaseURL + CBUserUrl

	authHeaders := cbAuth.NewAuthMap("GET", "", CBUserUrl)
	req, err := http.NewRequest("GET", url, nil)

	// Set auth headers
	for key, value := range authHeaders {
		log.Printf("Setting %s to %s", key, value)
		req.Header.Add(key, value)
	}

	// Retrieve response
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("%s", err)
		return "", ErrDecoding
	}
	defer resp.Body.Close()

	bodyAsStr, err := io.ReadAll(resp.Body)
	userResp := CBUserResponse{}

	if err := json.Unmarshal([]byte(bodyAsStr), &userResp); err != nil {
		log.Printf("error: %s", err)
		return "", ErrOnUnmarshall
	}

	log.Printf("Made it to the end")
	return userResp.Data.Id, nil
}

// CBRetrieveTransactions will return transactions for all coins the apikey has access to
func CBRetrieveTransactions(accountID string, cbAuth auth.CBAuth, client HttpClient) ([]CBTransaction, error) {

	transactionPath := strings.Replace(CBTransactionUrl, ":account_id", accountID, 0)

	url := CBBaseURL + transactionPath

	authHeaders := cbAuth.NewAuthMap("GET", "", CBUserUrl)
	req, err := http.NewRequest("GET", url, nil)

	// Set auth headers
	for key, value := range authHeaders {
		log.Printf("Setting %s to %s", key, value)
		req.Header.Add(key, value)
	}

	// Retrieve response
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("%s", err)
		return []CBTransaction{}, ErrDecoding
	}
	defer resp.Body.Close()

	bodyAsStr, err := io.ReadAll(resp.Body)
	var transactions []CBTransaction

	if err := json.Unmarshal([]byte(bodyAsStr), &transactions); err != nil {
		log.Printf("error: %s", err)
		return []CBTransaction{}, ErrOnUnmarshall
	}

	return transactions, nil
}

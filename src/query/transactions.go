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

// CBRetrieveTransactions will return transactions for all coins the apikey has access to
func CBRetrieveTransactions(accountID string, cbAuth auth.CBAuth, client HttpClient) ([]CBTransaction, error) {

	transactionPath := strings.Replace(CBTransactionUrl, ":account_id", accountID, -1)

	url := CBBaseURL + transactionPath

	authHeaders := cbAuth.NewAuthMap("GET", "", transactionPath)
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
		log.Printf("transaction_path: %s", transactionPath)
		log.Printf("url: %s", url)
		log.Printf("account_id: %s", accountID)
		log.Printf("error: %s", err)
		log.Printf("Body of response: %s", bodyAsStr)
		return []CBTransaction{}, ErrOnUnmarshall
	}

	return transactions, nil
}

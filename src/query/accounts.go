package query

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"warchest/src/auth"
)

// CBAccountsURL is the path to the GET accounts API call
const CBAccountsURL = "/v2/accounts"

// CBAccountsResp is the response object returned by retreiving all accounts
// TODO: need to handle pagination if there are that many accounts
type CBAccountsResp struct {
	Pagination struct {
		EndingBefore         interface{} `json:"ending_before"`
		StartingAfter        interface{} `json:"starting_after"`
		PreviousEndingBefore interface{} `json:"previous_ending_before"`
		NextStartingAfter    string      `json:"next_starting_after"`
		Limit                int         `json:"limit"`
		Order                string      `json:"order"`
		PreviousURI          interface{} `json:"previous_uri"`
		NextURI              string      `json:"next_uri"`
	} `json:"pagination"`
	Accounts []CBAccount `json:"data"`
}

// CBAccount is an individual account object
type CBAccount struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Primary  bool   `json:"primary"`
	Type     string `json:"type"`
	Currency struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Color        string `json:"color"`
		SortIndex    int    `json:"sort_index"`
		Exponent     int    `json:"exponent"`
		Type         string `json:"type"`
		AddressRegex string `json:"address_regex"`
		AssetID      string `json:"asset_id"`
		Slug         string `json:"slug"`
	} `json:"currency"`
	Balance struct {
		Amount   float64 `json:"amount,string"`
		Currency string  `json:"currency"`
	} `json:"balance"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Resource         string    `json:"resource"`
	ResourcePath     string    `json:"resource_path"`
	AllowDeposits    bool      `json:"allow_deposits"`
	AllowWithdrawals bool      `json:"allow_withdrawals"`
}

// CBRetrieveAccounts is a function that will retrieve all accounts associated with the account
func CBRetrieveAccounts(cbAuth auth.CBAuth, client HTTPClient) (CBAccountsResp, error) {

	url := CBBaseURL + CBAccountsURL
	authHeaders := cbAuth.NewAuthMap("GET", "", CBAccountsURL)
	req, err := http.NewRequest("GET", url, nil)

	// Set auth headers
	for key, value := range authHeaders {
		log.Printf("Setting %s to %s", key, value)
		req.Header.Add(key, value)
	}

	// per https://developers.coinbase.com/api/v2?shell#versioning
	t := time.Now()
	utcDate := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	req.Header.Add("CB-VERSION", utcDate)

	// Retrieve response
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Hit error on retrieval: %s", err)
		return CBAccountsResp{}, ErrConnection
	}
	defer resp.Body.Close()

	cResp := CBAccountsResp{}
	bodyAsStr, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body of error: %s", err)
		return CBAccountsResp{}, ErrDecoding
	}

	if err := json.Unmarshal([]byte(bodyAsStr), &cResp); err != nil {
		log.Printf("Body as string: \n%s\n", bodyAsStr)
		log.Printf("Couldn't unmarshall: %s", err)
		return CBAccountsResp{}, ErrOnUnmarshall
	}

	return cResp, err
}

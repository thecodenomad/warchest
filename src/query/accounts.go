package query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	"warchest/src/auth"
)

const CBAccountsUrl = "/v2/accounts"

type Accounts struct {
	Pagination struct {
		EndingBefore  interface{} `json:"ending_before"`
		StartingAfter interface{} `json:"starting_after"`
		Limit         int         `json:"limit"`
		Order         string      `json:"order"`
		PreviousUri   interface{} `json:"previous_uri"`
		NextUri       interface{} `json:"next_uri"`
	} `json:"pagination"`
	Data []struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Primary  bool   `json:"primary"`
		Type     string `json:"type"`
		Currency string `json:"currency"`
		Balance  struct {
			Amount   float64 `json:"amount,string"`
			Currency string  `json:"currency"`
		} `json:"balance"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Resource     string    `json:"resource"`
		ResourcePath string    `json:"resource_path"`
		Ready        bool      `json:"ready,omitempty"`
	} `json:"data"`
}

func CBRetrieveAccounts(cbAuth auth.CBAuth, client HttpClient) (Accounts, error) {

	url := CBBaseURL + CBAccountsUrl
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
		log.Printf("Hit error on retrieval: %s", err)
		return Accounts{}, ErrConnection
	}
	defer resp.Body.Close()

	cResp := Accounts{}

	bodyAsStr, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body of error: %s", err)
		return Accounts{}, ErrDecoding
	}

	if err := json.Unmarshal([]byte(bodyAsStr), &cResp); err != nil {
		log.Printf("Body as string: \n%s\n", bodyAsStr)
		log.Printf("Couldn't unmarshall: %s", err)
		return Accounts{}, ErrOnUnmarshall
	}

	return cResp, err
}

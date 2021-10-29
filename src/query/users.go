package query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"warchest/src/auth"
)

// CBRetrieveUserID will return the userID associated with the given api key
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

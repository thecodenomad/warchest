package query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"warchest/src/auth"
)

const CBUserUrl = "/v2/user"

type CBUserResp struct {
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
	userResp := CBUserResp{}

	if err := json.Unmarshal([]byte(bodyAsStr), &userResp); err != nil {
		log.Printf("error: %s", err)
		return "", ErrOnUnmarshall
	}

	return userResp.Data.Id, nil
}

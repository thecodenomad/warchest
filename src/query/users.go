package query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"warchest/src/auth"
)

// CBUserURL is the url path for retreiving a user object
const CBUserURL = "/v2/user"

// CBUserResp is the response object from a request to retrive a user object
type CBUserResp struct {
	Data struct {
		ID              string      `json:"id"`
		Name            string      `json:"name"`
		Username        string      `json:"username"`
		ProfileLocation interface{} `json:"profile_location"`
		ProfileBio      interface{} `json:"profile_bio"`
		ProfileURL      string      `json:"profile_url"`
		AvatarURL       string      `json:"avatar_url"`
		Resource        string      `json:"resource"`
		ResourcePath    string      `json:"resource_path"`
	} `json:"data"`
}

// CBRetrieveUserID will return the userID associated with the given api key
func CBRetrieveUserID(cbAuth auth.CBAuth, client HTTPClient) (string, error) {

	url := CBBaseURL + CBUserURL

	authHeaders := cbAuth.NewAuthMap("GET", "", CBUserURL)
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

	log.Printf("Body as string: %s", bodyAsStr)
	return userResp.Data.ID, nil
}

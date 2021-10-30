package query

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
	"time"
	"warchest/src/auth"
)

func TestCBRetrieveUserID(t *testing.T) {

	client := http.Client{
		Timeout: time.Second * 10,
	}
	var absClient HTTPClient
	absClient = &client

	t.Run("Happy Path", func(t *testing.T) {
		// Setup Test specifics
		cbAuth := auth.CBAuth{}
		json := `{ "data": { "id": "9da7a204-544e-5fd1-9a12-61176c5d4cd8" } }`

		// Estbalish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBUserURL).
			Reply(200).
			BodyString(json)

		actualResp, _ := CBRetrieveUserID(cbAuth, absClient)
		expectedResp := "9da7a204-544e-5fd1-9a12-61176c5d4cd8"
		assert.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cloudy Path", func(t *testing.T) {
		// Setup Test specifics
		cbAuth := auth.CBAuth{}
		json := `{unparseable,,,}`

		// Estbalish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBUserURL).
			Reply(200).
			BodyString(json)

		actualResp, _ := CBRetrieveUserID(cbAuth, absClient)
		expectedResp := ""
		assert.Equal(t, expectedResp, actualResp)
	})
}

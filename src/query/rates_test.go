package query

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
	"time"
)

// TODO: should be using a mock (spy) so we aren't making http requests
func TestRetrieveCoinData(t *testing.T) {
	// Test variables
	symbol := "ETH"

	client := http.Client{
		Timeout: time.Second * 10,
	}

	var absClient HTTPClient
	absClient = &client

	t.Run("Happy Path", func(t *testing.T) {

		json := `{"data": { "currency": "ETH", "rates": {"USD": "12.0", "EUR": "11.0", "GBP": "10.0"}}}`
		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBExchangeRateURL).
			Reply(200).
			BodyString(json)

		coinRates, err := CBRetrieveCoinRate(symbol, absClient)
		assert.Nil(t, err, "failed to retrieve rates")
		assert.Equal(t, 12.00, coinRates.USD, "Should be the same")
		assert.Equal(t, 11.00, coinRates.EUR, "Should be the same")
		assert.Equal(t, 10.00, coinRates.GBP, "Should be the same")
		assert.NotNil(t, coinRates, "no rates found!")
	})

	t.Run("Rainy Day connectivity!", func(t *testing.T) {

		mockClient := &MockClient{}

		_, err := CBRetrieveCoinRate(symbol, mockClient)

		// There should have been a connection error
		assert.Equal(t, ErrConnection, err, "this should be a connection error")
	})

	t.Run("Malformed JSON response", func(t *testing.T) {

		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBExchangeRateURL).
			Reply(200).
			BodyString(`[asdf,[],!}`)

		_, err := CBRetrieveCoinRate(symbol, absClient)

		// There should have been a connection error
		assert.Equal(t, ErrOnUnmarshall, err, "This call should have produced a JSON parse error")
	})

	// TODO: Still having problems forcing failure with io.ReadAll need to look into this
	t.Run("Failure to read response body check", func(t *testing.T) {

		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBExchangeRateURL).
			Reply(200).
			SetHeader("Content-Length", "1")

		_, err := CBRetrieveCoinRate(symbol, absClient)

		assert.NotNil(t, err, "This call should have produced a read error for the response body")
	})

	// addmittingly overkill, and borderline useful, but it makes for full coverage!
	t.Run("Test error handling", func(t *testing.T) {
		valueTests := []struct {
			actualValue   string
			expectedValue string
		}{
			{ErrDecoding.Error(), "failed decoding response"},
			{ErrOnUnmarshall.Error(), "failed to unmarshall"},
			{ErrConnection.Error(), "error during request"},
		}
		// Validate the rest of the imported values
		for _, tt := range valueTests {
			assert.Equal(t, tt.expectedValue, tt.actualValue)
		}
	})
}

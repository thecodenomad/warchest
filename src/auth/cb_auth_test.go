package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestCBAuth(t *testing.T) {

	// Setup Test specifics
	requestBody := ""
	requestMethod := "GET"
	requestPath := "/something?silly=true&currency=BTC"
	testAPIKey := "SoMeThInGcRaZy"
	testSecretKey := "aReYoUnOtEnTeRtAiNeD"

	// 30 second threshold, otherwise considered expired
	timestamp := int(time.Now().Unix())
	actualResp := CBAuth(testAPIKey, testSecretKey, requestMethod, requestBody, requestPath, timestamp)

	t.Run("Test return contents", func(t *testing.T) {
		responseTimestamp, _ := strconv.Atoi(actualResp[cbaccess_timestamp])

		// Make sure the required headers exist
		assert.Contains(t, actualResp, cbaccess_key)
		assert.Contains(t, actualResp, cbaccess_sign)
		assert.Contains(t, actualResp, cbaccess_timestamp)
		assert.Equal(t, timestamp, responseTimestamp, "timestamps should be the same!")

	})

	t.Run("Validate signature was calculated correctly", func(t *testing.T) {
		// Setup decoder ring - TODO: seems like there should be a better way to test this...
		sigText := "%d" + requestMethod + requestPath + requestBody
		h := hmac.New(sha256.New, []byte(testSecretKey))
		h.Write([]byte(sigText))
		expectedSignature := h.Sum(nil)
		actualSignature, _ := hex.DecodeString(actualResp[cbaccess_sign])

		assert.Equal(t, expectedSignature, actualSignature, "signatures should be the same")
	})
}

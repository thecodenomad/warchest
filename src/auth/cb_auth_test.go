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

	auth := CBAuth{testAPIKey, testSecretKey}

	// Generate NewAuthMap
	timestamp := int(time.Now().Unix())
	actualResp := auth.NewAuthMap(requestMethod, requestBody, requestPath)

	t.Run("Test return contents exist", func(t *testing.T) {
		// Make sure the required headers exist
		assert.Contains(t, actualResp, CBAccessKey)
		assert.Contains(t, actualResp, CBAccessSign)
		assert.Contains(t, actualResp, CBAccessTimestamp)
	})

	t.Run("Validate signature was calculated correctly", func(t *testing.T) {
		// Setup decoder ring - TODO: seems like there should be a better way to test this...
		sigText := strconv.Itoa(timestamp) + requestMethod + requestPath + requestBody
		h := hmac.New(sha256.New, []byte(testSecretKey))
		h.Write([]byte(sigText))
		expectedSignature := h.Sum(nil)
		actualSignature, _ := hex.DecodeString(actualResp[CBAccessSign])

		assert.Equal(t, expectedSignature, actualSignature, "signatures should be the same")
	})
}

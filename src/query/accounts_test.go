package query

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
	"time"
	auth2 "warchest/src/auth"
)

func TestCBRetrieveAccounts(t *testing.T) {

	cbAuth := auth2.CBAuth{"TestKey", "TestSecret"}
	client := http.Client{
		Timeout: time.Second * 10,
	}

	var absClient HttpClient
	absClient = &client

	t.Run("Happy Path", func(t *testing.T) {
		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBAccountsUrl).
			Reply(200).
			BodyString(validJson)

		accountResp, err := CBRetrieveAccounts(cbAuth, absClient)

		assert.Nil(t, err, "This is a happy path test, this shouldn't have had an err")
		assert.Equal(t, "58542935-67b5-56e1-a3f9-42686e07fa40", accountResp.Accounts[0].Id, "should be the same")
		assert.Equal(t, "My Vault", accountResp.Accounts[0].Name, "should be the same")
		assert.Equal(t, false, accountResp.Accounts[0].Primary, "should be the same")
		assert.Equal(t, "vault", accountResp.Accounts[0].Type, "should be the same")
		assert.Equal(t, "CTSI", accountResp.Accounts[0].Currency.Code, "should be the same")
		assert.Equal(t, 0.00000000, accountResp.Accounts[0].Balance.Amount, "should be the same")
		assert.Equal(t, "CTSI", accountResp.Accounts[0].Balance.Currency, "should be the same")
	})

	t.Run("Rainy Day connectivity!", func(t *testing.T) {

		mockClient := &MockClient{}

		accounts, err := CBRetrieveAccounts(cbAuth, mockClient)

		// There should have been a connection error
		assert.Equal(t, CBAccountsResp{}, accounts)
		assert.Equal(t, ErrConnection, err, "this should be a connection error")
	})

	t.Run("Malformed JSON response", func(t *testing.T) {

		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBAccountsUrl).
			Reply(200).
			BodyString(`[asdf,[],!}`)

		accounts, err := CBRetrieveAccounts(cbAuth, absClient)

		// There should have been a connection error
		assert.Equal(t, CBAccountsResp{}, accounts)
		assert.Equal(t, ErrOnUnmarshall, err, "This call should have produced a JSON parse error")
	})

	// TODO: Still having problems forcing failure with io.ReadAll need to look into this
	t.Run("Failure to read response body check", func(t *testing.T) {

		// Establish Mock
		defer gock.Off()
		gock.New(CBBaseURL).
			Get(CBAccountsUrl).
			Reply(200).
			SetHeader("Content-Length", "1")

		accounts, err := CBRetrieveAccounts(cbAuth, absClient)

		assert.NotNil(t, err, "This call should have produced a read error for the response body")
		assert.Equal(t, CBAccountsResp{}, accounts)
	})
}

const validJson = `{"pagination":{"ending_before":null,"starting_after":null,"previous_ending_before":null,"next_starting_after":"8e52c96f-bda7-508a-8876-972158d5e5c4","limit":25,"order":"desc","previous_uri":null,"next_uri":"/v2/accounts?starting_after=8e52c96f-bda7-508a-8876-972158d5e5c4"},"data":[{"id":"58542935-67b5-56e1-a3f9-42686e07fa40","name":"My Vault","primary":false,"type":"vault","currency":{"code":"CTSI","name":"Cartesi","color":"#1A1B1D","sort_index":173,"exponent":8,"type":"crypto","address_regex":"^(?:0x)?[0-9a-fA-F]{40}$","asset_id":"test-some-thing-else","slug":"cartesi"},"balance":{"amount":"0.00000000","currency":"CTSI"},"created_at":"2021-10-20T23:34:36Z","updated_at":"2021-10-27T01:47:08Z","resource":"account","resource_path":"/v2/accounts/some-thing-really-long-and-annoying","allow_deposits":true,"allow_withdrawals":true}]}`

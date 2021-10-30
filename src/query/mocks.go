package query

import (
	"errors"
	"net/http"
)

// MockClient helps mock failures at the transport level
type MockClient struct{}

// Do is the mocked method that forces a connection error
func (m *MockClient) Do(_ *http.Request) (*http.Response, error) {
	return nil, errors.New("test Connection Error")
}

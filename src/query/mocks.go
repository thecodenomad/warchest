package query

import (
	"errors"
	"net/http"
)

type MockClient struct{}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("test Connection Error")
}

package consumer

import (
	"net/http"
)

// Used to mock out HTTP requests to avoid external calls in unit tests
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func createConsumerWithMockedClient() (*Consumer, error) {
	c, err := NewConsumer()
	if err != nil {
		return nil, err
	}

	// used a mocked HTTP client
	c.Client = &MockClient{}
	return c, nil
}

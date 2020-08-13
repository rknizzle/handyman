package consumer

import (
	"net/http"
	"testing"
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

// test that sendRequestToApp calls client.Do
func TestSendRequestToAppCallsClientDoUnit(t *testing.T) {
	c, err := createConsumerWithMockedClient()
	if err != nil {
		t.Fatal(err)
	}

	// mock the Do function to verify that its called
	called := false
	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		called = true
		return nil, nil
	}

	_, err = c.sendRequestToApp("")
	if err != nil {
		t.Fatal(err)
	}

	if called != true {
		t.Fatal("Client.Do was not correctly called")
	}
}

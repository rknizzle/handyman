package consumer

import (
	"bytes"
	"errors"
	"io/ioutil"
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

// test that sendRequestToApp correctly passing the task message into the HTTP
// request body
func TestSendRequestToAppUsesTaskMessageUnit(t *testing.T) {
	c, err := createConsumerWithMockedClient()
	if err != nil {
		t.Fatal(err)
	}

	tests := []string{
		"{\"hello\":\"world!\"}",
		"",
	}

	for i, tc := range tests {
		t.Run(string(i), func(t *testing.T) {

			// mock the Do function to verify that it is passed the task message correctly
			GetDoFunc = func(req *http.Request) (*http.Response, error) {
				// get the request body
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				// verify that the request passed into client.Do contains the task message
				// in the body
				if string(body) != tc {
					t.Fatalf("Expected %s in request body but got %s", tc, string(body))
				}
				return nil, nil
			}

			_, err = c.sendRequestToApp(tc)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestHandleTaskUnit(t *testing.T) {
	c, err := createConsumerWithMockedClient()
	if err != nil {
		t.Fatal(err)
	}

	type mockedHTTPResponse struct {
		Resp *http.Response
		Err  error
	}

	type test struct {
		name         string
		httpResponse mockedHTTPResponse
		wantErr      bool
	}

	tests := []test{
		test{
			name: "200",
			httpResponse: mockedHTTPResponse{
				&http.Response{
					Body:   ioutil.NopCloser(bytes.NewReader([]byte{})),
					Status: "200"},
				nil,
			},
			wantErr: false,
		},
		test{
			name: "201",
			httpResponse: mockedHTTPResponse{
				&http.Response{
					Body:   ioutil.NopCloser(bytes.NewReader([]byte{})),
					Status: "201",
				},
				nil,
			},
			wantErr: false,
		},
		test{
			name:         "400",
			httpResponse: mockedHTTPResponse{&http.Response{Status: "400"}, nil},
			wantErr:      true,
		},
		test{
			name:         "500",
			httpResponse: mockedHTTPResponse{&http.Response{Status: "500"}, nil},
			wantErr:      true,
		},
		test{
			name:         "error",
			httpResponse: mockedHTTPResponse{nil, errors.New("")},
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// mock the HTTP requests return value to avoid external calls in unit test
			GetDoFunc = func(req *http.Request) (*http.Response, error) {
				return tc.httpResponse.Resp, tc.httpResponse.Err
			}

			_, gotErr := c.handleTask("")

			// test that handleTask returns an error when an error is expected (an
			// error is expected when the HTTP request returns a non 2xx status code)
			if tc.wantErr == (gotErr == nil) {
				t.Fatalf("%s: Wanted an error: %v. Got error value of: %v", tc.name, tc.wantErr, gotErr)
			}
		})
	}
}

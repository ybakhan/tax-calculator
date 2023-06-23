// Package taxclient provides functions to access resources from interview server
package taxclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitializeTaxClient(t *testing.T) {
	logger := log.NewNopLogger()
	httpClient := retryablehttp.NewClient()
	taxClient := InitializeTaxClient("http://interview-test-server:5000", httpClient, logger)
	assert.NotNil(t, taxClient)
}

func TestGetBrackets(t *testing.T) {
	logger := log.NewNopLogger()
	tests := map[string]struct {
		HTTPResponse     *http.Response
		HTTPError        error
		ReturnsError     bool
		ExpectedResponse GetBracketsResponse
		ExpectedBrackets []Bracket
	}{
		"get brackets failed": {
			HTTPResponse:     &http.Response{},
			HTTPError:        errors.New("some error"),
			ExpectedResponse: Failed,
			ReturnsError:     true,
		},
		"bracket not found": {
			HTTPResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			ExpectedResponse: NotFound,
		},
		"get brackets failed - server response not ok": {
			HTTPResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("some error")),
			},
			ExpectedResponse: Failed,
		},
		"get brackets failed - invalid json response": {
			HTTPResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("invalid json")),
			},
			ExpectedResponse: Failed,
			ReturnsError:     true,
		},
		"get brackets": {
			HTTPResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{\"tax_brackets\":[{\"max\":100392,\"min\":50197,\"rate\":0.205}]}")),
			},
			ExpectedResponse: Found,
			ExpectedBrackets: []Bracket{{Min: 50197, Max: 100392, Rate: 0.205}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockHTTPClient := &mockRetryableHTTPClient{}
			mockHTTPClient.
				On("Do", mock.AnythingOfType("*retryablehttp.Request")).
				Return(test.HTTPResponse, test.HTTPError)

			client := InitializeTaxClient("http://interview-test-server:5000", mockHTTPClient, logger)
			taxBrackets, response, err := client.GetBrackets(context.Background(), "2022")
			assert.Equal(t, test.ExpectedBrackets, taxBrackets)
			assert.Equal(t, test.ExpectedResponse, response)

			if test.ReturnsError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			mockHTTPClient.AssertExpectations(t)
		})
	}
}

type mockRetryableHTTPClient struct {
	mock.Mock
}

func (m *mockRetryableHTTPClient) Do(req *retryablehttp.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

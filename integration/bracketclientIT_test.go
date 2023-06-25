//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

func TestGetBrackets(t *testing.T) {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryWaitMin = 1 * time.Second
	httpClient.RetryWaitMax = 5 * time.Second
	httpClient.RetryMax = 5
	httpClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp.StatusCode == http.StatusInternalServerError {
			return true, nil
		}
		return false, nil
	}

	logger := log.NewNopLogger()
	bracketClient := taxbracket.InitializeBracketClient(os.Getenv("INTERVIEW_SERVER"), httpClient, logger)

	tests := map[string]struct {
		Year             string
		ExpectedResponse taxbracket.GetBracketsResponse
		ExpectedBrackets int
	}{
		"tax bracket not found": {
			Year:             "2018",
			ExpectedResponse: taxbracket.NotFound,
		},
		"get brackets": {
			Year:             "2022",
			ExpectedResponse: taxbracket.Found,
			ExpectedBrackets: 5,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			taxBrackets, response, err := bracketClient.GetBrackets(context.Background(), test.Year)
			assert.Equal(t, test.ExpectedBrackets, len(taxBrackets))
			assert.Equal(t, test.ExpectedResponse, response)
			assert.Nil(t, err)
		})
	}
}

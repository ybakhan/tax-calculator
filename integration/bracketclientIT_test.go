//go:build integration
// +build integration

package integration

import (
	"context"
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
	httpClient.HTTPClient.Timeout = 5000 * time.Millisecond
	httpClient.RetryWaitMin = 1000 * time.Millisecond
	httpClient.RetryWaitMax = 5000 * time.Millisecond
	httpClient.RetryMax = 5

	logger := log.NewNopLogger()
	bracketClient := taxbracket.InitializeBracketClient(os.Getenv("INTERVIEW_SERVER"), httpClient.StandardClient(), logger)

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

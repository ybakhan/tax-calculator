package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

func TestInitializeBracketCache(t *testing.T) {
	logger := log.NewNopLogger()
	getHandler := func(context.Context, string) (string, GetBracketsResponse) {
		return "", Found
	}
	saveHandler := func(context.Context, string, interface{}) error {
		return nil
	}
	bracketClient := InitializeBracketCache(getHandler, saveHandler, logger)
	assert.NotNil(t, bracketClient)
}

func TestGet(t *testing.T) {
	logger := log.NewNopLogger()
	ctx := context.Background()
	year := "2022"
	brackets := "[{\"min\":0,\"max\":50197,\"rate\":0.15},{\"min\":50197,\"max\":100392,\"rate\":0.205}]"

	tests := map[string]struct {
		Brackets         string
		ExpectedResponse GetBracketsResponse
		ExpectedBrackets []taxbracket.Bracket
	}{
		"Not found": {brackets, NotFound, nil},
		"Error":     {brackets, Failed, nil},
		"Found": {
			brackets,
			Found,
			[]taxbracket.Bracket{
				{
					Min:  0,
					Max:  50197,
					Rate: 0.15,
				},
				{
					Min:  50197,
					Max:  100392,
					Rate: 0.205,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			getHandler := func(context.Context, string) (string, GetBracketsResponse) {
				return test.Brackets, test.ExpectedResponse
			}
			bracketClient := InitializeBracketCache(getHandler, nil, logger)
			brackets, resp := bracketClient.Get(ctx, year)
			assert.Equal(t, test.ExpectedBrackets, brackets)
			assert.Equal(t, test.ExpectedResponse, resp)
		})
	}
}

func TestSave(t *testing.T) {
	logger := log.NewNopLogger()
	ctx := context.Background()
	year := "2022"

	brackets := []taxbracket.Bracket{
		{
			Min:  0,
			Max:  50197,
			Rate: 0.15,
		},
		{
			Min:  50197,
			Max:  100392,
			Rate: 0.205,
		},
	}

	tests := map[string]struct {
		Brackets         []taxbracket.Bracket
		ExpectedResponse SaveBracketsResponse
		ExpectedError    error
	}{
		"Save error": {
			brackets,
			SaveError,
			errors.New("some error"),
		},
		"Empty brackets": {
			nil,
			NotSaved,
			nil,
		},
		"Saved": {
			brackets,
			Saved,
			nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			saveHandler := func(context.Context, string, interface{}) error {
				return test.ExpectedError
			}
			bracketClient := InitializeBracketCache(nil, saveHandler, logger)
			resp, err := bracketClient.Save(ctx, year, test.Brackets)
			assert.Equal(t, test.ExpectedResponse, resp)
			assert.Equal(t, test.ExpectedError, err)
		})
	}

}

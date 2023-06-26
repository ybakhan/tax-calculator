package cache

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/log"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

func InitializeBracketCache(getHandler GetHandler, saveHandler SaveHandler, logger log.Logger) BracketCache {
	return &bracketCache{getHandler, saveHandler, logger}
}

// Get retrieves a tax bracket from cache for a given year
func (c *bracketCache) Get(ctx context.Context, year string) ([]taxbracket.Bracket, GetBracketsResponse) {
	bracketsStr, resp := c.GetHandler(ctx, year)
	if resp != Found {
		return nil, resp
	}

	var taxbrackets []taxbracket.Bracket
	err := json.Unmarshal([]byte(bracketsStr), &taxbrackets)
	if err != nil {
		c.Logger.Log("requestID", getRequestID(ctx), "error", err, "message", "error getting tax brackets from cache", "year", year)
		return nil, Failed
	}

	c.Logger.Log("requestID", getRequestID(ctx), "message", "tax brackets retrieved from cache", "taxbrackets", taxbrackets)
	return taxbrackets, Found
}

// Save saves tax brackets in cache for a given year
func (c *bracketCache) Save(ctx context.Context, year string, brackets []taxbracket.Bracket) (resp SaveBracketsResponse, err error) {
	defer func() {
		if err != nil {
			c.Logger.Log("requestID", getRequestID(ctx), "error", err, "message", "error saving tax brackets to cache", "year", year, "taxbrackets", brackets)
		}
	}()

	if len(brackets) == 0 {
		c.Logger.Log("requestID", getRequestID(ctx), "message", "empty tax brackets not saved in cache", "year", year, "taxbrackets", brackets)
		return NotSaved, nil
	}

	jsonBytes, err := json.Marshal(brackets)
	if err != nil {
		return NotSaved, err
	}

	err = c.SaveHandler(ctx, year, jsonBytes)
	if err != nil {
		return SaveError, err
	}

	c.Logger.Log("requestID", getRequestID(ctx), "message", "tax brackets saved in cache", "year", year, "taxbrackets", brackets)
	return Saved, nil
}

func getRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value("requestID").(string)
	if !ok {
		return ""
	}
	return requestID
}

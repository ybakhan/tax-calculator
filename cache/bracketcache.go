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

func (c *bracketCache) Get(ctx context.Context, year string) []taxbracket.Bracket {
	bracketsStr, resp := c.GetHandler(ctx, year)
	if resp != Found {
		return nil
	}

	var taxbrackets []taxbracket.Bracket
	err := json.Unmarshal([]byte(bracketsStr), &taxbrackets)
	if err != nil {
		c.Logger.Log("error", err, "message", "error getting tax brackets from cache", "year", year)
		return nil
	}

	c.Logger.Log("message", "tax brackets retrieved from cache", "taxbrackets", taxbrackets)
	return taxbrackets
}

func (c *bracketCache) Save(ctx context.Context, year string, brackets []taxbracket.Bracket) (err error) {
	defer func() {
		if err != nil {
			c.Logger.Log("error", err, "message", "error saving tax brackets to cache", "year", year, "taxbrackets", brackets)
		}
	}()

	jsonBytes, err := json.Marshal(brackets)
	if err != nil {
		return err
	}

	err = c.SaveHandler(ctx, year, jsonBytes)
	if err != nil {
		return err
	}

	c.Logger.Log("message", "tax brackets saved in cache", "year", year, "taxbrackets", brackets)
	return nil
}

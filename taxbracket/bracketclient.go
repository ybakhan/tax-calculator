package taxbracket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ybakhan/tax-calculator/common"
)

const bracketsResourcePath = "/tax-calculator/tax-year/"

func InitializeBracketClient(baseURL string, client retryableHTTPClient, logger log.Logger) BracketClient {
	bracketsURL, err := url.JoinPath(baseURL, bracketsResourcePath)
	if err != nil {
		err = fmt.Errorf("error intializing tax bracket client: %w", err)
		panic(err)
	}

	return &bracketClient{
		bracketsURL,
		client,
		logger,
	}
}

// GetBrackets gets tax brackets from the interview server
func (c *bracketClient) GetBrackets(ctx context.Context, year string) ([]Bracket, GetBracketsResponse, error) {
	brackets, response, err := c.getBrackets(ctx, year)
	if err != nil {
		c.logger.Log("requestID", common.GetRequestID(ctx), "error", err)
	}
	return brackets, response, err
}

func (c *bracketClient) getBrackets(ctx context.Context, year string) ([]Bracket, GetBracketsResponse, error) {
	taxBracketsURL, err := url.JoinPath(c.bracketsURL, year)
	if err != nil {
		return nil, Failed, err
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, taxBracketsURL, nil)
	if err != nil {
		return nil, Failed, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, Failed, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.logger.Log("requestID", common.GetRequestID(ctx), "msg", "tax brackets not found", "year", year)
		return nil, NotFound, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, Failed, nil
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, Failed, err
	}

	var taxbrackets Brackets
	err = json.Unmarshal(respBytes, &taxbrackets)
	if err != nil {
		return nil, Failed, err
	}

	if len(taxbrackets.Data) == 0 {
		c.logger.Log("requestID", common.GetRequestID(ctx), "msg", "tax brackets not found", "year", year)
		return nil, NotFound, nil
	}

	c.logger.Log("requestID", common.GetRequestID(ctx), "msg", "tax brackets found", "year", year, "taxBrackets", taxbrackets)
	return taxbrackets.Data, Found, nil
}

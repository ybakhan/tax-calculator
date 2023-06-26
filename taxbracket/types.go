package taxbracket

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-retryablehttp"
)

// GetBracketsResponse represents response type of get tax brackets function
type GetBracketsResponse int

const (
	Found GetBracketsResponse = -(iota)
	NotFound
	Failed
)

// BracketClient allows getting tax brackets for a given year
type BracketClient interface {
	GetBrackets(context.Context, string) ([]Bracket, GetBracketsResponse, error)
}

type retryableHTTPClient interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
}

type bracketClient struct {
	bracketsURL string
	client      retryableHTTPClient
	logger      log.Logger
}

// Bracket represents a tax bracket
type Bracket struct {
	Min  float64 `json:"min" example:"50197"`
	Max  float64 `json:"max" example:"100392"`
	Rate float64 `json:"rate" example:"0.205"`
}

type Brackets struct {
	Data []Bracket `json:"tax_brackets"`
}

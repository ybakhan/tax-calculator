package taxclient

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

type TaxClient interface {
	GetBrackets(context.Context, string) ([]Bracket, GetBracketsResponse, error)
}

type retryableHTTPClient interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
}

type taxClient struct {
	taxBracketsURL string
	client         retryableHTTPClient
	logger         log.Logger
}

// Bracket represents a tax bracket
type Bracket struct {
	Min  float32 `json:"min"`
	Max  float32 `json:"max"`
	Rate float32 `json:"rate"`
}

type Brackets struct {
	Data []Bracket `json:"tax_brackets"`
}

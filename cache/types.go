package cache

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

// GetBracketsResponse represents response type of get tax brackets from cache function
type GetBracketsResponse int

const (
	Found GetBracketsResponse = -(iota)
	NotFound
	GetError
)

// SaveBracketsResponse represents response type of save tax brackets to cache function
type SaveBracketsResponse int

const (
	Saved = -(iota)
	NotSaved
	SaveError
)

// BracketCache allows storage and retrieval of tax brackets from a cache
type BracketCache interface {
	Get(context.Context, string) ([]taxbracket.Bracket, GetBracketsResponse)
	Save(context.Context, string, []taxbracket.Bracket) (SaveBracketsResponse, error)
}

type bracketCache struct {
	GetHandler  GetHandler
	SaveHandler SaveHandler
	Logger      log.Logger
}

type GetHandler func(context.Context, string) (string, GetBracketsResponse)
type SaveHandler func(context.Context, string, interface{}) error

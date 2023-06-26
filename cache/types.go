package cache

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

type GetBracketsResponse int

const (
	Found GetBracketsResponse = -(iota)
	NotFound
	Failed
)

type SaveBracketsResponse int

const (
	Saved = -(iota)
	NotSaved
	SaveError
)

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

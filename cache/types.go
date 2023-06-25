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

type BracketCache interface {
	Get(context.Context, string) []taxbracket.Bracket
	Save(context.Context, string, []taxbracket.Bracket) error
}

type bracketCache struct {
	GetHandler  GetHandler
	SaveHandler SaveHandler
	Logger      log.Logger
}

type GetHandler func(context.Context, string) (string, GetBracketsResponse)
type SaveHandler func(context.Context, string, interface{}) error

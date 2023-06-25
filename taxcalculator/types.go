package taxcalculator

import (
	"github.com/ybakhan/tax-calculator/taxclient"
)

// TaxCalculation represents tax calculation for a year and salary
type TaxCalculation struct {
	TotalTaxes    float32      `json:"total_taxes" example:"8514.17"`
	EffectiveRate float32      `json:"effective_rate" example:"0.15"`
	BracketTaxes  []BracketTax `json:"taxes_by_band"`
}

// BracketTax represents tax calculated for a tax bracket
type BracketTax struct {
	Tax     float32           `json:"tax" example:"984.62"`
	Bracket taxclient.Bracket `json:"band"`
}

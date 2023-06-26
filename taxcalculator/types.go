package taxcalculator

import (
	"github.com/ybakhan/tax-calculator/taxbracket"
)

// TaxCalculation represents tax calculation for a year and salary
type TaxCalculation struct {
	Salary        float64      `json:"salary" example:"55000"`
	TotalTaxes    float64      `json:"total_taxes" example:"8514.17"`
	EffectiveRate float64      `json:"effective_rate" example:"0.15"`
	BracketTaxes  []BracketTax `json:"taxes_by_band,omitempty"`
}

// BracketTax represents tax calculated for a tax bracket
type BracketTax struct {
	Tax     float64            `json:"tax" example:"984.62"`
	Bracket taxbracket.Bracket `json:"band"`
}

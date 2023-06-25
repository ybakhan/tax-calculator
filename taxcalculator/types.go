package taxcalculator

import (
	"github.com/ybakhan/tax-calculator/taxclient"
)

// TaxCalculation represents tax calculation for a year and salary
type TaxCalculation struct {
	TotalTaxes    float32      `json:"total_taxes"`
	EffectiveRate float32      `json:"effective_rate"`
	BracketTaxes  []BracketTax `json:"taxes_by_band"`
}

// BracketTax represents tax calculated for a tax bracket
type BracketTax struct {
	Tax     float32           `json:"tax"`
	Bracket taxclient.Bracket `json:"band"`
}

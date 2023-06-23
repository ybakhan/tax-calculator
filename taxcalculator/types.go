package taxcalculator

import (
	"github.com/ybakhan/tax-calculator/taxclient"
)

// TaxCalculation represents tax calculation for a year and salary
type TaxCalculation struct {
	TotalTaxes    string       `json:"total_taxes"`
	EffectiveRate string       `json:"effective_rate"`
	BracketTaxes  []BracketTax `json:"taxes_by_band"`
}

// BracketTax represents tax calculated for a tax bracket
type BracketTax struct {
	Tax     string            `json:"tax"`
	Bracket taxclient.Bracket `json:"band"`
}

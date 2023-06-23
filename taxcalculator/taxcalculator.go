// Package taxcalculator provides functions for tax calculation
package taxcalculator

import (
	"fmt"
	"math"

	"github.com/ybakhan/tax-calculator/taxclient"
)

// Calculate computes taxes for a salary given tax brackets
func Calculate(brackets []taxclient.Bracket, salary float32) *TaxCalculation {
	if salary <= 0 {
		return nil
	}

	var taxByBand []BracketTax
	var total float32
	for _, bracket := range brackets {
		tax := calculateBracketTax(bracket, salary)
		if tax != 0 {
			taxByBand = append(taxByBand, BracketTax{
				format(tax),
				bracket,
			})
			total += tax
		}
	}

	answer := &TaxCalculation{
		format(total),
		format(total / salary),
		taxByBand,
	}
	return answer
}

func calculateBracketTax(bracket taxclient.Bracket, salary float32) float32 {
	if salary == 0 || salary <= bracket.Min {
		return 0
	}

	if salary > bracket.Min && (salary <= bracket.Max || bracket.Max == 0) {
		return round(bracket.Rate * (salary - bracket.Min))
	}

	return round(bracket.Rate * (bracket.Max - bracket.Min))
}

func round(f float32) float32 {
	return float32(math.Round(float64(f*100)) / 100)
}

func format(f float32) string {
	return fmt.Sprintf("%.2f", f)
}

// Package taxcalculator provides functions for tax calculation
package taxcalculator

import (
	"math"

	"github.com/ybakhan/tax-calculator/taxbracket"
)

// Calculate computes taxes for a salary given tax brackets
func Calculate(brackets []taxbracket.Bracket, salary float32) *TaxCalculation {
	if salary <= 0 {
		return nil
	}

	var taxByBand []BracketTax
	var total float32
	for _, bracket := range brackets {
		tax := calculateBracketTax(bracket, salary)
		if tax != 0 {
			taxByBand = append(taxByBand, BracketTax{
				tax,
				bracket,
			})
			total += tax
		}
	}

	answer := &TaxCalculation{
		salary,
		total,
		round(total / salary),
		taxByBand,
	}
	return answer
}

// calculateBracketTax calculates braket tax for a given salary
// result is rounded to 2 decimal places
func calculateBracketTax(bracket taxbracket.Bracket, salary float32) float32 {
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

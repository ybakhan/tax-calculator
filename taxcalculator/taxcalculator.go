// Package taxcalculator provides functions for tax calculation
package taxcalculator

import (
	"math"

	"github.com/ybakhan/tax-calculator/taxbracket"
)

// Calculate computes taxes for a salary given tax brackets
func Calculate(brackets []taxbracket.Bracket, salary float64) *TaxCalculation {
	if salary <= 0 {
		return &TaxCalculation{Salary: salary}
	}

	var taxByBand []BracketTax
	var total float64
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
		round(total),
		round(total / salary),
		taxByBand,
	}

	for i := range taxByBand {
		band := &taxByBand[i]
		band.Tax = round(band.Tax)
	}
	return answer
}

// calculateBracketTax calculates braket tax for a given salary
func calculateBracketTax(bracket taxbracket.Bracket, salary float64) float64 {
	if salary == 0 || salary <= bracket.Min {
		return 0
	}

	if salary > bracket.Min && (salary <= bracket.Max || bracket.Max == 0) {
		return bracket.Rate * (salary - bracket.Min)
	}

	return bracket.Rate * (bracket.Max - bracket.Min)
}

func round(f float64) float64 {
	return math.Round(f*100) / 100
}

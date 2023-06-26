package taxcalculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ybakhan/tax-calculator/taxbracket"
	"github.com/ybakhan/tax-calculator/testcommon"
)

func TestCalculate(t *testing.T) {
	taxBrackets := testcommon.ReadTaxBrackets(t, "../testcommon/taxbrackets.json")

	tests := map[string]struct {
		Salary   float64
		Expected TaxCalculation
	}{
		"calculate over one band": {
			50000,
			TaxCalculation{
				50000,
				7500,
				0.15,
				[]BracketTax{
					{7500, taxBrackets.Data[0]},
				},
			},
		},
		"calculate over one band with boundary salary": {
			50197,
			TaxCalculation{
				50197,
				7529.55,
				0.15,
				[]BracketTax{
					{7529.55, taxBrackets.Data[0]},
				},
			},
		},
		"calculate over two bands": {
			100000,
			TaxCalculation{
				100000,
				17739.17,
				0.18,
				[]BracketTax{
					{7529.55, taxBrackets.Data[0]},
					{10209.62, taxBrackets.Data[1]},
				},
			},
		},
		"calculate over two bands with boundary salary": {
			100392,
			TaxCalculation{
				100392,
				17819.52,
				0.18,
				[]BracketTax{
					{7529.55, taxBrackets.Data[0]},
					{10289.97, taxBrackets.Data[1]},
				},
			},
		},
		"calculate over three bands": {
			100393,
			TaxCalculation{
				100393,
				17819.78,
				0.18,
				[]BracketTax{
					{7529.55, taxBrackets.Data[0]},
					{10289.97, taxBrackets.Data[1]},
					{0.26, taxBrackets.Data[2]},
				},
			},
		},
		"calculate over five bands": {
			1234567,
			TaxCalculation{
				1234567,
				385587.65,
				0.31,
				[]BracketTax{
					{7529.55, taxBrackets.Data[0]},
					{10289.97, taxBrackets.Data[1]},
					{14360.58, taxBrackets.Data[2]},
					{19164.07, taxBrackets.Data[3]},
					{334243.47, taxBrackets.Data[4]},
				},
			},
		},
		"zero salary": {
			0, TaxCalculation{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			taxCalculation := Calculate(taxBrackets.Data, test.Salary)
			assert.Equal(t, &test.Expected, taxCalculation)
		})
	}
}

func TestCalculateByBracket(t *testing.T) {
	tests := map[string]struct {
		Bracket     taxbracket.Bracket
		Salary      float64
		ExpectedTax float64
	}{
		"first bracket": {
			Salary: 55000,
			Bracket: taxbracket.Bracket{
				Min:  0,
				Max:  50197,
				Rate: 0.15,
			},
			ExpectedTax: 7529.55,
		},
		"second bracket": {
			Salary: 55000,
			Bracket: taxbracket.Bracket{
				Min:  50197,
				Max:  100392,
				Rate: 0.205,
			},
			ExpectedTax: 984.61,
		},
		"out of bracket": {
			Salary: 50197,
			Bracket: taxbracket.Bracket{
				Min:  50197,
				Max:  100392,
				Rate: 0.205,
			},
			ExpectedTax: 0,
		},
		"bracket boundary": {
			Salary: 50197,
			Bracket: taxbracket.Bracket{
				Min:  0,
				Max:  50197,
				Rate: 0.15,
			},
			ExpectedTax: 7529.55,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tax := calculateBracketTax(test.Bracket, test.Salary)
			assert.Equal(t, test.ExpectedTax, round(tax))
		})
	}
}

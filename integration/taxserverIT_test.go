//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ybakhan/tax-calculator/taxcalculator"
	"github.com/ybakhan/tax-calculator/testcommon"
)

func TestGetTaxes(t *testing.T) {
	taxBrackets := testcommon.ReadTaxBrackets(t, "../testcommon/taxbrackets.json")

	tests := map[string]struct {
		salary   string
		Expected taxcalculator.TaxCalculation
	}{
		"calculate over one band": {
			"50196",
			taxcalculator.TaxCalculation{
				TotalTaxes:    "7529.40",
				EffectiveRate: "0.15",
				BracketTaxes: []taxcalculator.BracketTax{
					{
						Tax:     "7529.40",
						Bracket: taxBrackets.Data[0],
					},
				},
			},
		},
		"calculate over one band boundary": {
			"50197",
			taxcalculator.TaxCalculation{
				TotalTaxes:    "7529.55",
				EffectiveRate: "0.15",
				BracketTaxes: []taxcalculator.BracketTax{
					{Tax: "7529.55", Bracket: taxBrackets.Data[0]},
				},
			},
		},
		"calculate over two bands": {
			"55000",
			taxcalculator.TaxCalculation{
				TotalTaxes:    "8514.17",
				EffectiveRate: "0.15",
				BracketTaxes: []taxcalculator.BracketTax{
					{Tax: "7529.55", Bracket: taxBrackets.Data[0]},
					{Tax: "984.62", Bracket: taxBrackets.Data[1]},
				},
			},
		},
		"calculate over two bands boundary": {
			"100392",
			taxcalculator.TaxCalculation{
				TotalTaxes:    "17819.52",
				EffectiveRate: "0.18",
				BracketTaxes: []taxcalculator.BracketTax{
					{Tax: "7529.55", Bracket: taxBrackets.Data[0]},
					{Tax: "10289.97", Bracket: taxBrackets.Data[1]},
				},
			},
		},
		"calculate over three bands": {
			"100393",
			taxcalculator.TaxCalculation{
				TotalTaxes:    "17819.78",
				EffectiveRate: "0.18",
				BracketTaxes: []taxcalculator.BracketTax{
					{Tax: "7529.55", Bracket: taxBrackets.Data[0]},
					{Tax: "10289.97", Bracket: taxBrackets.Data[1]},
					{Tax: "0.26", Bracket: taxBrackets.Data[2]},
				},
			},
		},
		"calculate over five bands": {
			"221709",
			taxcalculator.TaxCalculation{
				TotalTaxes:    "51344.50",
				EffectiveRate: "0.23",
				BracketTaxes: []taxcalculator.BracketTax{
					{Tax: "7529.55", Bracket: taxBrackets.Data[0]},
					{Tax: "10289.97", Bracket: taxBrackets.Data[1]},
					{Tax: "14360.58", Bracket: taxBrackets.Data[2]},
					{Tax: "19164.07", Bracket: taxBrackets.Data[3]},
					{Tax: "0.33", Bracket: taxBrackets.Data[4]},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			body := getTaxes(t, "2022", test.salary)
			taxCalculation := taxcalculator.TaxCalculation{}
			err := json.Unmarshal(body, &taxCalculation)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, test.Expected, taxCalculation)
		})
	}
}

func TestGetTaxes_TaxYearNotFound(t *testing.T) {
	body := getTaxes(t, "2023", "50000")
	expected := "{\"message\":\"tax year not found 2023\"}\n"
	assert.Equal(t, expected, string(body))
}

func TestGetTaxes_TaxYearNotProvided(t *testing.T) {
	body := getTaxes(t, "", "50000")
	assert.Equal(t, "404 page not found\n", string(body))
}

func TestGetTaxes_InvalidTaxYear(t *testing.T) {
	body := getTaxes(t, "abc", "50000")
	expected := "{\"error\":\"invalid tax year abc\"}\n"
	assert.Equal(t, expected, string(body))
}

func TestGetTaxes_SalaryNotProvided(t *testing.T) {
	body := getTaxes(t, "2023", "")
	expected := "{\"error\":\"salary missing in request\"}\n"
	assert.Equal(t, expected, string(body))
}

func TestGetTaxes_SalaryInvalid(t *testing.T) {
	body := getTaxes(t, "2023", "abc")
	expected := "{\"error\":\"invalid salary abc\"}\n"
	assert.Equal(t, expected, string(body))
}

func getTaxes(t *testing.T, year, salary string) []byte {
	client := &http.Client{}
	params := url.Values{}
	params.Add("s", salary)

	requestURL := fmt.Sprintf("%s/tax/%s?%s", os.Getenv("TAX_CALCULATOR_SERVER"), year, params.Encode())
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

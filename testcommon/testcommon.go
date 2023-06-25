package testcommon

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/ybakhan/tax-calculator/taxbracket"
)

func ReadTaxBrackets(t *testing.T, path string) *taxbracket.Brackets {
	file, _ := ioutil.ReadFile(path)
	var taxBrackets taxbracket.Brackets
	err := json.Unmarshal([]byte(file), &taxBrackets)
	if err != nil {
		t.Fatal("Error parsing test input file")
	}
	return &taxBrackets
}

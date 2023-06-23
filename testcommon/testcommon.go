package testcommon

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/ybakhan/tax-calculator/taxclient"
)

func ReadTaxBrackets(t *testing.T, path string) *taxclient.Brackets {
	file, _ := ioutil.ReadFile(path)
	var taxBrackets taxclient.Brackets
	err := json.Unmarshal([]byte(file), &taxBrackets)
	if err != nil {
		t.Fatal("Error parsing test input file")
	}
	return &taxBrackets
}

package main

import (
	"context"
	"encoding/json"

	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ybakhan/tax-calculator/cache"
	"github.com/ybakhan/tax-calculator/taxbracket"
	"github.com/ybakhan/tax-calculator/taxcalculator"
	"github.com/ybakhan/tax-calculator/testcommon"
)

func TestHandleGetTaxes(t *testing.T) {
	taxBrackets := testcommon.ReadTaxBrackets(t, "../testcommon/taxbrackets.json")
	salaryStr := "1234567"
	expectedTaxes := &taxcalculator.TaxCalculation{
		Salary:        1234567,
		TotalTaxes:    385587.65,
		EffectiveRate: 0.31,
		BracketTaxes: []taxcalculator.BracketTax{
			{7529.55, taxBrackets.Data[0]},
			{10289.97, taxBrackets.Data[1]},
			{14360.58, taxBrackets.Data[2]},
			{19164.07, taxBrackets.Data[3]},
			{334243.47, taxBrackets.Data[4]},
		},
	}

	logger := log.NewNopLogger()
	tests := map[string]struct {
		Year                         string
		Salary                       string
		CachedBrackets               []taxbracket.Bracket
		Brackets                     []taxbracket.Bracket
		GetBracketsFromCacheResponse cache.GetBracketsResponse
		GetBracketsResponse          taxbracket.GetBracketsResponse
		GetBracketsError             error
		ExpectedTaxes                *taxcalculator.TaxCalculation
		ExpectedStatusCode           int
		ExpectedResponse             string
	}{
		"invalid year": {
			Year:               "abc",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   "{\"error\":\"invalid tax year abc\"}\n",
		},
		"missing salary": {
			Year:               "2022",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   "{\"error\":\"salary missing in request\"}\n",
		},
		"invalid salary": {
			Year:               "2022",
			Salary:             "abc",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   "{\"error\":\"invalid salary abc\"}\n",
		},
		"brackets found in cache": {
			Year:                         "2022",
			Salary:                       salaryStr,
			CachedBrackets:               taxBrackets.Data,
			GetBracketsFromCacheResponse: cache.Found,
			ExpectedStatusCode:           http.StatusOK,
			ExpectedTaxes:                expectedTaxes,
		},
		"brackets not found in cache, get brackets success": {
			Year:                         "2022",
			Salary:                       salaryStr,
			GetBracketsFromCacheResponse: cache.NotFound,
			Brackets:                     taxBrackets.Data,
			GetBracketsResponse:          taxbracket.Found,
			ExpectedStatusCode:           http.StatusOK,
			ExpectedTaxes:                expectedTaxes,
		},
		"brackets not found in cache, get brackets error": {
			Year:                         "2022",
			Salary:                       salaryStr,
			GetBracketsFromCacheResponse: cache.NotFound,
			GetBracketsError:             errors.New("some error"),
			ExpectedStatusCode:           http.StatusInternalServerError,
			ExpectedResponse:             "{\"error\":\"some error\"}\n",
		},
		"brackets not found in cache, get brackets failed": {
			Year:                         "2022",
			Salary:                       salaryStr,
			GetBracketsFromCacheResponse: cache.NotFound,
			GetBracketsResponse:          taxbracket.Failed,
			ExpectedStatusCode:           http.StatusInternalServerError,
			ExpectedResponse:             "{\"error\":\"get taxes failed year 2022\"}\n",
		},
		"brackets not found in cache, tax year not found": {
			Year:                         "2023",
			Salary:                       salaryStr,
			GetBracketsFromCacheResponse: cache.NotFound,
			GetBracketsResponse:          taxbracket.NotFound,
			ExpectedStatusCode:           http.StatusNotFound,
			ExpectedResponse:             "{\"message\":\"tax year not found 2023\"}\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockBracketClient := mockBracketClient{}
			mockBracketCache := mockBracketCache{}

			if test.ExpectedStatusCode != http.StatusBadRequest {
				mockBracketCache.
					On("Get", mock.Anything, test.Year).
					Return(test.CachedBrackets, test.GetBracketsFromCacheResponse)

				if len(test.CachedBrackets) == 0 {
					mockBracketClient.
						On("GetBrackets", mock.Anything, test.Year).
						Return(test.Brackets, test.GetBracketsResponse, test.GetBracketsError)

					if test.ExpectedStatusCode != http.StatusInternalServerError &&
						test.ExpectedStatusCode != http.StatusNotFound {
						mockBracketCache.
							On("Save", mock.Anything, test.Year, test.Brackets).
							Return(cache.Saved, nil)
					}
				}
			}

			s := &taxServer{"", &mockBracketClient, &mockBracketCache, nil, logger}

			router := mux.NewRouter()
			router.HandleFunc("/tax/{year}", s.makeHTTPHandlerFunc(s.handleGetTaxes))

			request, _ := http.NewRequest("GET", fmt.Sprintf("/tax/%s?s=%s", test.Year, test.Salary), nil)
			vars := map[string]string{
				"year": test.Year,
			}
			request = mux.SetURLVars(request, vars)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			result := recorder.Result()
			assert.Equal(t, test.ExpectedStatusCode, result.StatusCode)

			body, _ := ioutil.ReadAll(result.Body)
			if test.ExpectedTaxes != nil {
				var taxes taxcalculator.TaxCalculation
				json.Unmarshal(body, &taxes)
				assert.Equal(t, *test.ExpectedTaxes, taxes)
			} else {
				assert.Equal(t, test.ExpectedResponse, string(body))
			}

			mockBracketClient.AssertExpectations(t)
			mockBracketCache.AssertExpectations(t)
		})
	}
}

type mockBracketClient struct {
	mock.Mock
}

func (c *mockBracketClient) GetBrackets(ctx context.Context, year string) ([]taxbracket.Bracket, taxbracket.GetBracketsResponse, error) {
	args := c.Called(ctx, year)
	return args.Get(0).([]taxbracket.Bracket), args.Get(1).(taxbracket.GetBracketsResponse), args.Error(2)
}

type mockBracketCache struct {
	mock.Mock
}

func (c *mockBracketCache) Get(ctx context.Context, year string) ([]taxbracket.Bracket, cache.GetBracketsResponse) {
	args := c.Called(ctx, year)
	return args.Get(0).([]taxbracket.Bracket), args.Get(1).(cache.GetBracketsResponse)
}

func (c *mockBracketCache) Save(ctx context.Context, year string, brackets []taxbracket.Bracket) (cache.SaveBracketsResponse, error) {
	args := c.Called(ctx, year, brackets)
	return cache.SaveBracketsResponse(args.Int(0)), args.Error(1)
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/ybakhan/tax-calculator/cache"
	"github.com/ybakhan/tax-calculator/common"
	_ "github.com/ybakhan/tax-calculator/docs"
	"github.com/ybakhan/tax-calculator/taxbracket"
	"github.com/ybakhan/tax-calculator/taxcalculator"
)

// Start starts rest api server that handles tax calculation requests
func (s *taxServer) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/login", s.makeHTTPHandlerFunc(s.handleLogin))
	router.HandleFunc("/tax/{year}", s.makeHTTPHandlerFunc(s.validateApiKey(s.handleGetTaxes)))
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	s.Logger.Log("msg", fmt.Sprintf("tax calculator listening on port %s", s.ListenAddress))
	return http.ListenAndServe(s.ListenAddress, router)
}

// handleGetTaxes handles get taxes api call go doc
//
//	@Summary		calculate taxes
//	@Description	calculate taxes for given a salary and tax year
//	@Tags			taxes
//	@Produce		json
//	@Param			year	path		int	true	"tax year"
//	@Param			s		query		int	true	"salary"
//	@Success		200		{object}	taxcalculator.TaxCalculation
//	@Failure		404		{object}	taxServerResponse
//	@Failure		500		{object}	taxServerError
//	@Router			/tax/{year} [get]
func (s *taxServer) handleGetTaxes(w http.ResponseWriter, r *http.Request) (resp *handlerResponse, err error) {
	if r.Method != "GET" {
		return &handlerResponse{Status: http.StatusMethodNotAllowed}, fmt.Errorf("method not supported %s", r.Method)
	}

	vars := mux.Vars(r)
	year := vars["year"]
	if _, err := strconv.Atoi(year); err != nil {
		return &handlerResponse{Status: http.StatusBadRequest}, fmt.Errorf("invalid tax year %s", year)
	}

	salaryStr := r.FormValue("s")
	if salaryStr == "" {
		return &handlerResponse{Status: http.StatusBadRequest}, errors.New("salary missing in request")
	}

	salaryF, err := strconv.ParseFloat(salaryStr, 64)
	if err != nil {
		return &handlerResponse{Status: http.StatusBadRequest}, fmt.Errorf("invalid salary %s", salaryStr)
	}

	ctx := r.Context()
	var taxes *taxcalculator.TaxCalculation
	if brackets, resp := s.BracketCache.Get(ctx, year); resp == cache.Found {
		taxes = taxcalculator.Calculate(brackets, salaryF)
	} else {
		brackets, response, err := s.BracketClient.GetBrackets(ctx, year)
		if err != nil {
			return &handlerResponse{Status: http.StatusInternalServerError}, err
		}

		if response == taxbracket.Failed {
			return &handlerResponse{Status: http.StatusInternalServerError}, fmt.Errorf("get taxes failed year %s", year)
		}

		if response == taxbracket.NotFound {
			notFoundMessage := fmt.Sprintf("tax year not found %s", year)
			s.Logger.Log("requestID", common.GetRequestID(ctx), "msg", notFoundMessage)
			return &handlerResponse{http.StatusNotFound, &taxServerResponse{notFoundMessage}}, nil
		}

		s.BracketCache.Save(ctx, year, brackets)
		taxes = taxcalculator.Calculate(brackets, salaryF)
	}

	s.Logger.Log("requestID", common.GetRequestID(ctx), "msg", "calculated taxes", "year", year, "salary", salaryF, "taxes", taxes)
	return &handlerResponse{http.StatusOK, taxes}, nil
}

// writeJSON sets status header and
// writes a reponse to a http response writer
func (s *taxServer) makeHTTPHandlerFunc(f requestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		s.Logger.Log("requestID", common.GetRequestID(ctx), "msg", "handling request", "method", r.Method, "url", r.URL.Path)

		resp, err := f(w, r.WithContext(ctx))
		var responseBody any
		if err != nil {
			s.Logger.Log("requestID", common.GetRequestID(ctx), "error", err)
			responseBody = &taxServerError{Error: err.Error()}
		} else {
			responseBody = resp.Body
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			s.Logger.Log("requestID", common.GetRequestID(ctx), "error", err)
		}
	}
}

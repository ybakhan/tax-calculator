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
	_ "github.com/ybakhan/tax-calculator/docs"
	"github.com/ybakhan/tax-calculator/taxbracket"
	"github.com/ybakhan/tax-calculator/taxcalculator"
)

// Start starts rest api server that handles tax calculation requests
func (s *taxServer) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/tax/{year}", s.makeHTTPHandlerFunc(s.handleTaxes))
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	s.Logger.Log("msg", fmt.Sprintf("tax calculator listening on port %s", s.ListenAddress))
	http.ListenAndServe(s.ListenAddress, router)
}

// handleTaxes handles api calls on /tax/{year}
func (s *taxServer) handleTaxes(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "GET" {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "requestID", requestID)

		s.Logger.Log("requestID", getRequestID(ctx), "msg", "handling request", "method", r.Method, "url", r.URL.Path)
		return s.handleGetTaxes(w, r.WithContext(ctx))
	}

	return http.StatusBadRequest, fmt.Errorf("method not supported %s", r.Method)
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
func (s *taxServer) handleGetTaxes(w http.ResponseWriter, r *http.Request) (i int, err error) {
	defer func() {
		if err != nil {
			s.Logger.Log("requestID", getRequestID(r.Context()), "error", err)
		}
	}()

	vars := mux.Vars(r)
	year := vars["year"]
	if _, err := strconv.Atoi(year); err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid tax year %s", year)
	}

	salaryStr := r.FormValue("s")
	if salaryStr == "" {
		return http.StatusBadRequest, errors.New("salary missing in request")
	}

	salaryF, err := strconv.ParseFloat(salaryStr, 64)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid salary %s", salaryStr)
	}

	ctx := r.Context()
	var taxes *taxcalculator.TaxCalculation
	if brackets, resp := s.BracketCache.Get(ctx, year); resp == cache.Found {
		taxes = taxcalculator.Calculate(brackets, salaryF)
	} else {
		brackets, response, err := s.BracketClient.GetBrackets(ctx, year)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if response == taxbracket.Failed {
			return http.StatusInternalServerError, fmt.Errorf("get taxes failed year %s", year)
		}

		if response == taxbracket.NotFound {
			notFoundMessage := fmt.Sprintf("tax year not found %s", year)
			s.Logger.Log("requestID", getRequestID(ctx), "msg", notFoundMessage)

			writeJSON(w, http.StatusNotFound, taxServerResponse{notFoundMessage})
			return http.StatusNotFound, nil
		}

		s.BracketCache.Save(ctx, year, brackets)
		taxes = taxcalculator.Calculate(brackets, salaryF)
	}

	s.Logger.Log("requestID", getRequestID(ctx), "msg", "calculated taxes", "year", year, "salary", salaryF, "taxes", taxes)
	writeJSON(w, http.StatusOK, taxes)
	return http.StatusOK, nil
}

func (s *taxServer) makeHTTPHandlerFunc(f requestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if status, err := f(w, r); err != nil {
			writeJSON(w, status, &taxServerError{Error: err.Error()})
		}
	}
}

// writeJSON sets status header and
// writes a reponse to a http response writer
func writeJSON(w http.ResponseWriter, status int, response any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}

func getRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value("requestID").(string)
	if !ok {
		return ""
	}
	return requestID
}

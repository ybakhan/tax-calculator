package main

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/ybakhan/tax-calculator/cache"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

// taxServer represents api server that handles requests to calculate taxes
type taxServer struct {
	ListenAddress string
	BracketClient taxbracket.BracketClient
	BracketCache  cache.BracketCache
	Logger        log.Logger
}

// taxServerError represents error response of tax server api
type taxServerError struct {
	Error string `json:"error"`
}

// taxServerResponse represents plain text response of tax server api
type taxServerResponse struct {
	Message string `json:"message"`
}

type requestHandler func(http.ResponseWriter, *http.Request) (int, error)

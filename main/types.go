package main

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

// taxServer represents api server that handles requests to calculate taxes
type taxServer struct {
	ListenAddress string
	TaxClient     taxbracket.BracketClient
	Redis         *redis.Client
	Logger        log.Logger
}

type taxServerError struct {
	Error string `json:"error"`
}

type taxServerResponse struct {
	Message string `json:"message"`
}

type requestHandler func(http.ResponseWriter, *http.Request) (int, error)

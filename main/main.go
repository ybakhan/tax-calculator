package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

//	@title			Tax Calculator API
//	@version		1.0
//	@description	REST API for calculating taxes

//	@contact.name	Yasser Khan
//	@contact.url	http://github.com/ybakhan
//	@contact.email	ybakhan@gmail.com
func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	config := readConfig()
	logger.Log("msg", "tax calculator started", "configuration", &config)

	httpClient := initializeHTTPClient(config)
	taxClient := taxbracket.InitializeBracketClient(config.InterviewServer.BaseURL, httpClient, logger)

	listenAddress := fmt.Sprintf(":%d", config.Port)
	server := &taxServer{listenAddress, taxClient, logger}
	server.Start()
}

func initializeHTTPClient(config *Config) *retryablehttp.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Timeout = time.Duration(config.HTTPClient.TimeoutMs) * time.Millisecond
	httpClient.RetryWaitMin = time.Duration(config.HTTPClient.Retry.Wait.MinMs) * time.Millisecond
	httpClient.RetryWaitMax = time.Duration(config.HTTPClient.Retry.Wait.MaxMs) * time.Millisecond
	httpClient.RetryMax = config.HTTPClient.Retry.Max
	return httpClient
}

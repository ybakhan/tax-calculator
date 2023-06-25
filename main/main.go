package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ybakhan/tax-calculator/taxclient"
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
	taxClient := taxclient.InitializeTaxClient(config.InterviewServer.BaseURL, httpClient, logger)

	listenAddress := fmt.Sprintf(":%d", config.Port)
	server := &taxServer{listenAddress, taxClient, logger}
	server.Start()
}

func initializeHTTPClient(config *Config) *retryablehttp.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryWaitMin = time.Duration(config.HTTPClient.Retry.Wait.MinSeconds) * time.Second
	httpClient.RetryWaitMax = time.Duration(config.HTTPClient.Retry.Wait.MaxSeconds) * time.Second
	httpClient.RetryMax = config.HTTPClient.Retry.Max

	httpClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp.StatusCode == http.StatusInternalServerError {
			return true, nil
		}
		return false, nil
	}
	return httpClient
}

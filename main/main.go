package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
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
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)

	config := readConfig()
	logger := initializeLogger()
	logger.Log("msg", "tax calculator started", "configuration", &config)

	redis := initializeRedis(config, logger)

	// Create a wait group to wait for the cleanup code to finish
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Wait for the terminate signal
		<-terminate

		// disconnect redis
		err := redis.Close()
		if err != nil {
			logger.Log("error", err, "msg", "error closing Redis connection")
			return
		}
		logger.Log("msg", "Redis connection closed")
	}()

	initializeTaxServer(config, redis, logger)

	wg.Wait()
	logger.Log("msg", "terminating tax-calculator")
}

func initializeTaxServer(config *Config, redis *redis.Client, logger log.Logger) {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Timeout = time.Duration(config.HTTPClient.TimeoutMs) * time.Millisecond
	httpClient.RetryWaitMin = time.Duration(config.HTTPClient.Retry.Wait.MinMs) * time.Millisecond
	httpClient.RetryWaitMax = time.Duration(config.HTTPClient.Retry.Wait.MaxMs) * time.Millisecond
	httpClient.RetryMax = config.HTTPClient.Retry.Max

	taxClient := taxbracket.InitializeBracketClient(config.InterviewServer.BaseURL, httpClient, logger)

	listenAddress := fmt.Sprintf(":%d", config.Port)
	server := &taxServer{listenAddress, taxClient, redis, logger}
	server.Start()
}

func initializeLogger() log.Logger {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}

func initializeRedis(config *Config, logger log.Logger) *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       0,
	})

	pong, err := redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	logger.Log("msg", "Connected to Redis", "pong", pong)
	return redis
}

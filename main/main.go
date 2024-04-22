package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ybakhan/tax-calculator/cache"
	"github.com/ybakhan/tax-calculator/common"
	"github.com/ybakhan/tax-calculator/taxbracket"
)

//	@title			Tax Calculator API
//	@version		1.0
//	@description	REST API for calculating taxes

// @contact.name	Yasser Khan
// @contact.url	http://github.com/ybakhan
// @contact.email	ybakhan@gmail.com
func main() {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)

	config := readConfig()
	logger := initializeLogger()
	logger.Log("msg", "tax calculator started", "configuration", &config)

	redis := initializeRedis(config, logger)

	go func() {
		<-terminate
		// disconnect redis
		if err := redis.Close(); err != nil {
			logger.Log("error", err, "msg", "error closing Redis connection")
			return
		}
		logger.Log("msg", "Redis connection closed")
	}()

	if err := initializeTaxServer(config, redis, logger); err != nil {
		if err != http.ErrServerClosed {
			logger.Log("error", err, "msg", "server encountered an error")
		} else {
			logger.Log("msg", "terminating tax-calculator")
			os.Exit(1)
		}
	}
}

func initializeTaxServer(config *Config, redisClient *redis.Client, logger log.Logger) error {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Timeout = time.Duration(config.HTTPClient.TimeoutMs) * time.Millisecond
	httpClient.RetryWaitMin = time.Duration(config.HTTPClient.Retry.Wait.MinMs) * time.Millisecond
	httpClient.RetryWaitMax = time.Duration(config.HTTPClient.Retry.Wait.MaxMs) * time.Millisecond
	httpClient.RetryMax = config.HTTPClient.Retry.Max

	listenAddress := fmt.Sprintf(":%d", config.Port)
	bracketClient := taxbracket.InitializeBracketClient(config.InterviewServer.BaseURL, httpClient.StandardClient(), logger)
	bracketCache := initializeBracketCache(redisClient, logger)

	server := &taxServer{listenAddress, bracketClient, bracketCache, &config.ApiToken, logger}
	return server.Start()
}

func initializeBracketCache(redisClient *redis.Client, logger log.Logger) cache.BracketCache {
	getHandler := func(ctx context.Context, year string) (string, cache.GetBracketsResponse) {
		result, err := redisClient.Get(ctx, year).Result()
		if err == redis.Nil {
			return "", cache.NotFound
		}

		if err != nil {
			logger.Log("requestID", common.GetRequestID(ctx), "error", err, "msg", "error getting tax brackets from cache")
			return "", cache.GetError
		}

		return result, cache.Found
	}

	saveHandler := func(ctx context.Context, year string, value interface{}) error {
		// tax brackets are cached indefinitely
		return redisClient.Set(ctx, year, value, 0).Err()
	}

	return cache.InitializeBracketCache(getHandler, saveHandler, logger)
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

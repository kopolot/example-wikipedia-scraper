package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"example-wikipedia-scraper/internal/cache"
	types "example-wikipedia-scraper/internal/types/api"
)

const (
	idempotencyKeyPrefix  = "idempotency:"
	idempotencyProcessing = "__processing__"
)

type cachedApiResponse struct {
	StatusCode int                    `json:"status_code"`
	Body       *types.ApiResponseBody `json:"body"`
	Headers    map[string]string      `json:"headers"`
}

func IdempotentMiddleware(request *types.ApiRequest) *types.ApiResponse {
	token, ok := getIdempotencyToken(request.Headers)
	if !ok || token == "" {
		return badRequestResponse("X-Idempotent-Token header is required")
	}

	redis := cache.GetRedisClient()
	if redis == nil {
		return internalErrorResponse()
	}

	ctx := context.Background()
	key := idempotencyKeyPrefix + token

	if cached, err := getCachedResponse(ctx, redis, key); err == nil && cached != nil {
		return cached
	}

	acquired, err := redis.SetNX(ctx, key, idempotencyProcessing, 5*time.Minute)
	if err != nil {
		return internalErrorResponse()
	}
	if !acquired {
		if cached, err := getCachedResponse(ctx, redis, key); err == nil && cached != nil {
			return cached
		}
		return conflictResponse("Request with this idempotency token is already being processed")
	}

	return nil
}

func CacheIdempotentResponse(request *types.ApiRequest, response *types.ApiResponse) {
	success := false
	token, ok := getIdempotencyToken(request.Headers)
	if !ok || token == "" || response == nil {
		return
	}

	key := idempotencyKeyPrefix + token

	redis := cache.GetRedisClient()
	if redis == nil {
		return
	}

	defer func() {
		if !success {
			redis.Del(context.Background(), key)
		}
	}()

	if !response.Body.Success {
		return
	}

	cached := cachedApiResponse{
		StatusCode: response.StatusCode,
		Body:       response.Body,
		Headers:    response.Headers,
	}
	data, err := json.Marshal(cached)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = redis.Set(ctx, key, string(data), 0)
	if err == nil {
		success = true
	}
}

func getIdempotencyToken(headers map[string]string) (string, bool) {
	for key, value := range headers {
		if strings.EqualFold(key, "X-Idempotent-Token") {
			return strings.TrimSpace(value), true
		}
	}
	return "", false
}

func getCachedResponse(ctx context.Context, redis *cache.RedisClient, key string) (*types.ApiResponse, error) {
	val, err := redis.Get(ctx, key)
	if err != nil || val == "" || val == idempotencyProcessing {
		return nil, err
	}

	var cached cachedApiResponse
	if err := json.Unmarshal([]byte(val), &cached); err != nil {
		return nil, err
	}

	return &types.ApiResponse{
		StatusCode: cached.StatusCode,
		Body:       cached.Body,
		Headers:    cached.Headers,
	}, nil
}

func badRequestResponse(msg string) *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusBadRequest,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: msg}},
		},
	}
}

func conflictResponse(msg string) *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusConflict,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: msg}},
		},
	}
}

func internalErrorResponse() *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusInternalServerError,
		Body: &types.ApiResponseBody{
			Success: false,
			Errors:  []types.ErrorDetail{{Message: "Internal Server Error"}},
		},
	}
}

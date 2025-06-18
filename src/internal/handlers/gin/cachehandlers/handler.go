/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 03 Apr 2024
*/

package cachehandlers

import (
	"github.com/gin-gonic/gin"
	reqmodels "lruCache/poc/src/internal/api/gin/models"
	"lruCache/poc/src/internal/helpers"
	"lruCache/poc/src/internal/web"
	cache2 "lruCache/poc/src/modules/cache"
	"strconv"
	"time"
)

func GetValue(ctx *gin.Context) (*web.JSONResponse, *web.ErrorStruct) {
	log := helpers.GetLogger()
	requestBody := web.GetJsonBody[reqmodels.GetModel](ctx)
	log.Debug(ctx, requestBody)
	cache := cache2.GetLRUCache(1024)

	log.Debug(ctx, "Cache", cache)

	value, found := cache.Get(requestBody.Key)
	if found {
		return web.NewHTTPSuccessResponse(ctx, "found matching key", map[string]interface{}{
			"key": value,
		}), nil
	}

	return web.NewHTTPSuccessResponse(ctx, "matching key not found", map[string]interface{}{
		"key": value,
	}), nil
}

func SetValue(ctx *gin.Context) (*web.JSONResponse, *web.ErrorStruct) {
	log := helpers.GetLogger()
	requestBody := web.GetJsonBody[reqmodels.SetModel](ctx)
	log.Debug(ctx, requestBody)
	cache := cache2.GetLRUCache(1024)

	log.Debug(ctx, "Cache", cache)

	duration, err := strconv.ParseInt(requestBody.Expiration, 10, 64)
	if err != nil {
		return web.NewHTTPFailureResponse(ctx, "error converting expiration from string to int64", err), nil
	}

	cache.Set(requestBody.Key, requestBody.Value, time.Duration(duration*1000000000))

	return web.NewHTTPSuccessResponse(ctx, "matching key not found", map[string]interface{}{
		requestBody.Key: requestBody.Value,
		"expiration":    requestBody.Expiration,
	}), nil
}

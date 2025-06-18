/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package gin

import (
	"fmt"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	"lruCache/poc/src/internal/helpers"
	"lruCache/poc/src/internal/web"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(context *gin.Context) (*web.JSONResponse, *web.ErrorStruct)

func ServeEndpoint[BodyType any](nextHandler func(ctx *gin.Context) (*web.JSONResponse, *web.ErrorStruct)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		log := helpers.GetLogger()
		if ctx.GetHeader(string(constants.XCorrelationID)) == "" {
			ctx.Set(string(constants.CorrelationID), helpers.GetTraceIdFromContext(ctx))
		} else {
			ctx.Set(string(constants.CorrelationID), ctx.GetHeader(string(constants.XCorrelationID)))
		}
		ctx.Set(string(constants.RequestID), ctx.GetHeader(string(constants.XRequestID)))
		ctx.Set(constants.AppName, fmt.Sprintf("%s:%s", config.AppName, constants.AppNameSuffixHTTPServer))
		ctx.Set(constants.AppVersion, ctx.GetHeader(string(constants.XAppVersion)))
		baseRoutePath := strings.Replace(ctx.FullPath(), config.BaseRouterPath, "", -1)
		splits := strings.SplitN(baseRoutePath, "/", 4)
		routeGroup := splits[2]
		routePath := splits[3]
		ctx.Set(constants.RouteGroup, routeGroup)
		ctx.Set(constants.RoutePath, routePath)
		if config.AppEnv != constants.EnvProd {
			log.Info(ctx, "request_received", map[string]interface{}{
				"method":   ctx.Request.Method,
				"uri":      ctx.Request.RequestURI,
				"clientIp": ctx.ClientIP(),
				"headers":  ctx.Request.Header,
				"params":   ctx.Params,
			})
		} else {
			log.Info(ctx, "request_received", map[string]interface{}{
				"method":   ctx.Request.Method,
				"uri":      ctx.Request.RequestURI,
				"clientIp": ctx.ClientIP(),
				"params":   ctx.Params,
			})
		}
		var response *web.JSONResponse
		var responseErr *web.ErrorStruct
		err := web.ValidateJsonBody[BodyType](ctx)
		if err != nil {
			responseErr = err
		} else {
			response, responseErr = nextHandler(ctx)
		}
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		if config.AppEnv != constants.EnvProd {
			code, respBody := web.BuildResponse(ctx, response, responseErr)
			log.Info(ctx, "request_processed", map[string]interface{}{
				"response": map[string]interface{}{
					"statusCode": code,
					"jsonBody":   respBody,
				},
				"latencyInMilli": latency.Milliseconds(),
				"method":         ctx.Request.Method,
				"uri":            ctx.Request.RequestURI,
			})
		} else {
			log.Info(ctx, "request_processed", map[string]interface{}{
				"latencyInMilli": latency.Milliseconds(),
				"method":         ctx.Request.Method,
				"uri":            ctx.Request.RequestURI,
			})
		}
		ctx.Header(fmt.Sprintf(string(constants.XTraceID)), fmt.Sprintf("x2mp-%s", helpers.GetTraceIdFromContext(ctx)))
		web.SendResponse(ctx, response, responseErr)
	}
}

/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package kafka

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	AppName        string = "appName"
	CorrelationID  string = "correlationId"
	XCorrelationID string = "x-correlation-id"
	RequestID      string = "requestId"
	XRequestID     string = "x-request-id"
)

func GetValueFromContext(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}

func GetDefaultValueFromContext(ctx context.Context, contextKey string) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetString(string(contextKey))
	} else {
		return GetValueFromContext(ctx, string(contextKey))
	}
}

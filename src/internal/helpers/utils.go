/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package helpers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"lruCache/poc/src/constants"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

func GenerateId(prefix string) string {
	uuidV4 := uuid.New().String()
	uuidWithoutHyphens := strings.ToUpper(strings.Replace(uuidV4, "-", "", -1))
	return fmt.Sprintf("%s%s", prefix, uuidWithoutHyphens)
}

func GetSHA256HashFromString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func GetDefaultValueFromContext(ctx context.Context, contextKey string) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetString(string(contextKey))
	} else {
		return getValueFromContext(ctx, string(contextKey))
	}
}

func GetTraceIdFromContext(ctx context.Context) string {
	var spanCtx = trace.SpanContextFromContext(ctx)
	if ginCtx, ok := ctx.(*gin.Context); ok {
		spanCtx = trace.SpanContextFromContext(ginCtx.Request.Context())
	}
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func GetSpanIdFromContext(ctx context.Context) string {
	var spanCtx = trace.SpanContextFromContext(ctx)
	if ginCtx, ok := ctx.(*gin.Context); ok {
		spanCtx = trace.SpanContextFromContext(ginCtx.Request.Context())
	}
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}

func getValueFromContext(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}

// IsExactlyOneTrue returns true if exactly one of the input booleans is true.
func IsExactlyOneTrue(a, b, c bool) bool {
	return (a && !b && !c) || (!a && b && !c) || (!a && !b && c)
}

func IsRoutePathPresentInRoutePathList(routePath constants.RoutesPath, routePathList []constants.RoutesPath) bool {
	for _, rp := range routePathList {
		if rp == routePath {
			return true
		}
	}
	return false
}

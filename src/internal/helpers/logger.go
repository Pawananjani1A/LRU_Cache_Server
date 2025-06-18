/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package helpers

import (
	"context"
	GLogger "lruCache/poc/lib/logger"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	"sync"
)

var loggerMutex = sync.RWMutex{}
var logger *GLogger.LoggerService

func GetLogger() *GLogger.LoggerService {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return logger
}

func setLogger(value *GLogger.LoggerService) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	logger = value
}

type MetaDataProvider struct {
}

func (l *MetaDataProvider) AppName(ctx context.Context) string {
	// Use a predefined value if not found in ctx
	appName := GetDefaultValueFromContext(ctx, constants.AppName)
	if appName == "" {
		return config.AppName
	}
	return appName
}
func (l *MetaDataProvider) UserID(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.UserID)
}
func (l *MetaDataProvider) CorrelationID(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.CorrelationID)
}
func (l *MetaDataProvider) LoggerName(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.LoggerName)
}
func (l *MetaDataProvider) TraceID(ctx context.Context) string {
	traceID := GetTraceIdFromContext(ctx)
	if traceID != "" {
		return traceID
	}
	return GetDefaultValueFromContext(ctx, constants.TraceID)
}
func (l *MetaDataProvider) SpanID(ctx context.Context) string {
	spanID := GetSpanIdFromContext(ctx)
	if spanID != "" {
		return spanID
	}
	return GetDefaultValueFromContext(ctx, constants.SpanID)
}
func (l *MetaDataProvider) Thread(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.Thread)
}
func (l *MetaDataProvider) ClientDeviceID(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.ClientDeviceID)
}
func (l *MetaDataProvider) ClientPlatform(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.ClientPlatform)
}
func (l *MetaDataProvider) ClientVersion(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.ClientVersion)
}
func (l *MetaDataProvider) AppPackageID(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.AppPackageID)
}
func (l *MetaDataProvider) ClientSessionID(ctx context.Context) string {
	return GetDefaultValueFromContext(ctx, constants.ClientSessionID)
}

func InitLogger() *GLogger.LoggerService {
	var metaDataProvider GLogger.ILoggingMetaDataProvider = &MetaDataProvider{}
	l := GLogger.NewLoggerService(config.LogLevel, metaDataProvider)
	setLogger(l)
	return l
}

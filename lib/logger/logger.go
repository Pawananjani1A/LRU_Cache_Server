/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package GLogger

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

// LogLevel represents the level of logging
type LogLevel int

// Enum for log levels
const (
	DEBUG LogLevel = iota + 1
	VERBOSE
	INFO
	WARN
	ERROR
)

// String method for LogLevel
func (l LogLevel) String() string {
	return [...]string{"DEBUG", "VERBOSE", "INFO", "WARN", "ERROR"}[l-1]
}

// ILoggingMetaDataProvider interface for building log metadata
type ILoggingMetaDataProvider interface {
	AppName(ctx context.Context) string
	UserID(ctx context.Context) string
	CorrelationID(ctx context.Context) string
	LoggerName(ctx context.Context) string
	TraceID(ctx context.Context) string
	SpanID(ctx context.Context) string
	Thread(ctx context.Context) string
	ClientDeviceID(ctx context.Context) string
	ClientPlatform(ctx context.Context) string
	ClientVersion(ctx context.Context) string
	AppPackageID(ctx context.Context) string
	ClientSessionID(ctx context.Context) string
}

// LoggerService struct
type LoggerService struct {
	metaDataProvider ILoggingMetaDataProvider
	appliedLogLevel  LogLevel
}

// NewLoggerService creates a new LoggerService
func NewLoggerService(logLevel LogLevel, metaDataProvider ILoggingMetaDataProvider) *LoggerService {
	return &LoggerService{
		metaDataProvider: metaDataProvider,
		appliedLogLevel:  logLevel,
	}
}

// buildFormattedMessage builds the formatted message
func (ls *LoggerService) buildFormattedMessage(message []interface{}) string {
	formattedMsg := ""
	for _, e := range message {
		switch v := e.(type) {
		case error:
			formattedMsg += fmt.Sprintf("%v \nStack Trace:\n %s \n", v.Error(), string(debug.Stack()))
			fmt.Println("error: ", v)
			fmt.Println(string(debug.Stack()))
		default:
			marshalV, _ := json.Marshal(v)
			formattedMsg += " " + string(marshalV)
		}
	}
	return strings.TrimSpace(formattedMsg)
}

// buildFormattedMetaData builds the formatted metadata
func (ls *LoggerService) buildFormattedMetaData(ctx context.Context, customMetadata map[string]interface{}, level LogLevel) map[string]interface{} {
	meta := map[string]interface{}{
		"timestamp":       time.Now().UnixNano() / int64(time.Millisecond),
		"application":     ls.metaDataProvider.AppName(ctx),
		"principal":       ls.metaDataProvider.UserID(ctx),
		"correlationId":   ls.metaDataProvider.CorrelationID(ctx),
		"loggerName":      ls.metaDataProvider.LoggerName(ctx),
		"traceId":         ls.metaDataProvider.TraceID(ctx),
		"spanId":          ls.metaDataProvider.SpanID(ctx),
		"clientDeviceId":  ls.metaDataProvider.ClientDeviceID(ctx),
		"clientPlatform":  ls.metaDataProvider.ClientPlatform(ctx),
		"clientVersion":   ls.metaDataProvider.ClientVersion(ctx),
		"appPackageId":    ls.metaDataProvider.AppPackageID(ctx),
		"clientSessionId": ls.metaDataProvider.ClientSessionID(ctx),
		"thread":          ls.metaDataProvider.Thread(ctx),
		"level":           level.String(),
	}
	for k, v := range customMetadata {
		meta[k] = v
	}
	return meta
}

// printLog prints the log
func (ls *LoggerService) printLog(ctx context.Context, additionalParams map[string]interface{}, level LogLevel, message ...interface{}) string {
	if level < ls.appliedLogLevel {
		return ""
	}
	formattedMessage := ls.buildFormattedMessage(message)
	metaData := ls.buildFormattedMetaData(ctx, additionalParams, level)
	metaData["message"] = formattedMessage
	formattedLog, _ := json.Marshal(metaData)
	finalLog := string(formattedLog)
	fmt.Println(finalLog)
	return finalLog
}

// Debug logs a debug message
func (ls *LoggerService) Debug(ctx context.Context, message ...interface{}) string {
	return ls.printLog(ctx, make(map[string]interface{}), DEBUG, message...)
}

// DebugWithContext logs a debug message with context
func (ls *LoggerService) DebugWithContext(ctx context.Context, customFields map[string]interface{}, message ...interface{}) string {
	return ls.printLog(ctx, customFields, DEBUG, message...)
}

// Verbose logs a verbose message
func (ls *LoggerService) Verbose(ctx context.Context, message ...interface{}) string {
	return ls.printLog(ctx, make(map[string]interface{}), INFO, message...)
}

// VerboseWithContext logs a verbose message with context
func (ls *LoggerService) VerboseWithContext(ctx context.Context, customFields map[string]interface{}, message ...interface{}) string {
	return ls.printLog(ctx, customFields, INFO, message...)
}

// Info logs an info message
func (ls *LoggerService) Info(ctx context.Context, message ...interface{}) string {
	return ls.printLog(ctx, make(map[string]interface{}), INFO, message...)
}

// InfoWithContext logs an info message with context
func (ls *LoggerService) InfoWithContext(ctx context.Context, customFields map[string]interface{}, message ...interface{}) string {
	return ls.printLog(ctx, customFields, INFO, message...)
}

// Warn logs a warning message
func (ls *LoggerService) Warn(ctx context.Context, message ...interface{}) string {
	return ls.printLog(ctx, make(map[string]interface{}), WARN, message...)
}

// WarnWithContext logs a warning message with context
func (ls *LoggerService) WarnWithContext(ctx context.Context, customFields map[string]interface{}, message ...interface{}) string {
	return ls.printLog(ctx, customFields, WARN, message...)
}

// Error logs an error message
func (ls *LoggerService) Error(ctx context.Context, message ...interface{}) string {
	return ls.printLog(ctx, make(map[string]interface{}), ERROR, message...)
}

// ErrorWithContext logs an error message with context
func (ls *LoggerService) ErrorWithContext(ctx context.Context, customFields map[string]interface{}, message ...interface{}) string {
	return ls.printLog(ctx, customFields, ERROR, message...)
}

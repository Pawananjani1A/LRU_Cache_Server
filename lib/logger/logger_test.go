/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package GLogger

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type MetaDataProvider struct {
}

func (l *MetaDataProvider) AppName(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) UserID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) CorrelationID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) LoggerName(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) TraceID(ctx context.Context) string {
	return uuid.NewString()
}
func (l *MetaDataProvider) SpanID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) Thread(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientDeviceID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientPlatform(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientVersion(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) AppPackageID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientSessionID(ctx context.Context) string {
	return ""
}

func getLogger(level LogLevel) *LoggerService {
	var metaDataProvider ILoggingMetaDataProvider = &MetaDataProvider{}
	return NewLoggerService(level, metaDataProvider)
}
func TestLoggerForMap(t *testing.T) {
	logger := getLogger(INFO)
	logM := logger.InfoWithContext(context.Background(), map[string]interface{}{
		"timestamp": 0,
	}, "hello", map[string]interface{}{
		"integer":  1,
		"name":     "pawananjani",
		"lastname": "kumar",
		"email":    "pawananjanimth1@gmail.com",
	})

	if logM != `{"appPackageId":"","application":"appname","clientDeviceId":"","clientPlatform":"","clientSessionId":"","clientVersion":"","correlationId":"correlationid","level":"INFO","loggerName":"","message":"\"hello\" {\"email\":\"pawananjanimth1@gmail.com\",\"integer\":1,\"lastname\":\"kumar\",\"name\":\"pawananjani\"}","principal":"userid","spanId":"","timestamp":0,"traceId":""}` {
		t.Errorf("log mismatch for logging for map")
	}
}

func TestLoggerForSliceString(t *testing.T) {
	logger := getLogger(INFO)
	logM := logger.InfoWithContext(context.Background(), map[string]interface{}{
		"timestamp": 0,
	}, "hello", []string{
		"pawananjani",
		"kumar",
		"1",
		"2",
	})
	if logM != `{"appPackageId":"","application":"appname","clientDeviceId":"","clientPlatform":"","clientSessionId":"","clientVersion":"","correlationId":"correlationid","level":"INFO","loggerName":"","message":"\"hello\" [\"pawananjani\",\"kumar\",\"1\",\"2\"]","principal":"userid","spanId":"","timestamp":0,"traceId":""}` {
		t.Errorf("log mismatch for logging for map")
	}
}

func TestLoggerForSliceMap(t *testing.T) {
	logger := getLogger(INFO)
	logM := logger.InfoWithContext(context.Background(), map[string]interface{}{
		"timestamp": 0,
	}, "hello", []map[string]interface{}{
		{
			"firstName": "pawananjani",
			"lastName":  "kumar",
			"email":     "pawananjanimth1@gmail.com",
		},
		{
			"firstName": "shivank",
			"lastName":  "singhal",
			"email":     "shivank.kumar@gmail.com",
		},
	})
	if logM != `{"appPackageId":"","application":"appname","clientDeviceId":"","clientPlatform":"","clientSessionId":"","clientVersion":"","correlationId":"correlationid","level":"INFO","loggerName":"","message":"\"hello\" [{\"email\":\"pawananjanimth1@gmail.com\",\"firstName\":\"pawananjani\",\"lastName\":\"kumar\"},{\"email\":\"shivank.kumar@gmail.com\",\"firstName\":\"shivank\",\"lastName\":\"singhal\"}]","principal":"userid","spanId":"","timestamp":0,"traceId":""}` {
		t.Errorf("log mismatch for logging for map")
	}
}

func TestLoggerService_Info(t *testing.T) {
	logger := getLogger(INFO)
	logger.Info(context.Background(), map[string]interface{}{
		"timestamp": 0,
	}, "hello", map[string]interface{}{
		"integer":  1,
		"name":     "pawananjani",
		"lastname": "kumar",
		"email":    "pawananjanimth1@gmail.com",
	})
}

func TestLoggerService_Debug(t *testing.T) {
	logger := getLogger(INFO)
	logM := logger.Debug(context.Background(), map[string]interface{}{
		"timestamp": 0,
	}, "hello", map[string]interface{}{
		"integer":  1,
		"name":     "pawananjani",
		"lastname": "kumar",
		"email":    "pawananjanimth1@gmail.com",
	})
	if logM != "" {
		t.Errorf("invalid: debug log should not be printed at info level")
	}
}

func TestLoggerService_ErrorWithContext(t *testing.T) {
	logger := getLogger(INFO)
	err := errors.New("testing an error string")
	logger.ErrorWithContext(context.Background(), map[string]interface{}{
		"timestamp": 0,
	}, "error_occurred", err)

}

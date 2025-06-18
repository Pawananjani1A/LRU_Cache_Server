/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

import (
	"context"
	GLogger "lruCache/poc/lib/logger"

	"github.com/google/uuid"
)

type MetaDataProvider struct{}

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
func getLogger(level GLogger.LogLevel) *GLogger.LoggerService {
	var metaDataProvider GLogger.ILoggingMetaDataProvider = &MetaDataProvider{}
	return GLogger.NewLoggerService(level, metaDataProvider)
}

var logger = getLogger(GLogger.DEBUG)

const (
	CommEventNameAllocateProductGroupSuccess CommunicationEventID = "COMM_EVENT_ALLOCATE_PRODUCT_GROUP_SUCCESS"
)

var testCommunicationTemplates = CommunicationTemplates{
	CommEventNameAllocateProductGroupSuccess: []CommunicationEventType{CommEventTypeSSE},
}

var testCommSSEEventNameMap = CommEventIDSSEEventMap{
	CommEventNameAllocateProductGroupSuccess: "sse_X2_UserModule_MultiProduct_AllocateProductGroup_Success",
}

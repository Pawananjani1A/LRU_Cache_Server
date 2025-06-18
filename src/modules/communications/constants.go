/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

const (
	CommEventTypeSSE         CommunicationEventType = "SSE"
	CommEventTypeGlobalAlert CommunicationEventType = "GLOBAL_ALERT"
	CommEventTypeAudit       CommunicationEventType = "AUDIT"
)

const (
	CommErrModuleNotInitiated            CommunicationError = "COMMUNICATIONS_MODULE_NOT_INITIATED"
	CommErrUnknownCommDestination        CommunicationError = "UNKNOWN_COMMUNICATION_DESTINATION_ENCOUNTERED"
	CommErrNoDestinationFound            CommunicationError = "NO_COMM_DESTINATION_FOUND"
	CommErrSSEEventBuilderMethodNotFound CommunicationError = "SSE_BUILDER_NOT_FOUND"
	CommErrSSEEventBuildFailed           CommunicationError = "SSE_EVENT_BUILD_FAILED"
)

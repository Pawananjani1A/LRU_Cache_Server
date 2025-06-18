/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

import "context"

type CommUserDetails struct {
	AppUserID string `json:"appUserId,omitempty"`
}
type CommunicationEventType string
type CommunicationEventID string
type CommunicationPayload map[string]interface{}
type CommunicationTemplates map[CommunicationEventID][]CommunicationEventType
type CommunicationError string

// communicationEventID to eventName mapping for each event type

type CommEventIDSSEEventMap map[CommunicationEventID]string
type CommSSEEventBuilder func(ctx context.Context, userDetails CommUserDetails, eventID CommunicationEventID, payload CommunicationPayload) (*CommSSEEventPayload, error)

type CommEventIDGlobalAlertMap map[CommunicationEventID]string
type CommGlobalAlertBuilder func()

type CommEventIDAuditLogMap map[CommunicationEventID]string
type CommAuditLogBuilder func()

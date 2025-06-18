/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

import (
	"context"
	"errors"
	"fmt"
	"lruCache/poc/src/constants"
	"lruCache/poc/src/internal/helpers"
	"lruCache/poc/src/modules/communications"
	"time"
)

func BuildSSEEvent(ctx context.Context, userDetails communications.CommUserDetails, commEventID communications.CommunicationEventID, commPayload communications.CommunicationPayload) (*communications.CommSSEEventPayload, error) {
	logger := helpers.GetLogger()
	sseEventName := constants.CommSSEEventNameMap[commEventID]
	if sseEventName == "" {
		logger.Error(ctx, fmt.Sprintf("no_sse_event_found_for_commEventID=%s", commEventID))
		return nil, errors.New("SSE_EVENT_NOT_REGISTERED_FOR_EVENT_ID___BuildSSEEvent")
	}
	eventPayload := communications.CommSSEEventPayload{
		PackageName:    "x2-core-multi-product",
		AppUserID:      userDetails.AppUserID,
		PayloadVersion: "1.0",
		Profile:        "",
		Events:         nil,
	}
	switch commEventID {
	case constants.CommEventNameAddedToGroupSuccess:
		eventPayload.Events = []communications.CommSSEEventData{
			{
				TimeStamp:      time.Now().UnixMilli(),
				EventName:      sseEventName,
				PayloadVersion: "1.0",
				EventData:      commPayload,
			},
		}
		break
	default:
		logger.Warn(ctx, "unexpected_event_id...unreachable_condition_reached...")
		return nil, errors.New("UNEXPECTED_EVENT_ID_RECEIVED___BuildSSEEvent")
	}
	return &eventPayload, nil
}

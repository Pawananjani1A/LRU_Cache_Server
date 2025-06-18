/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

import (
	"context"
	"errors"
	"fmt"
	GLogger "lruCache/poc/lib/logger"
)

type communicationsDispatcher struct {
	appUserID      string
	logger         *GLogger.LoggerService
	communications ICommunications
}

type ICommunicationsDispatcher interface {
	SendCommunicationEvent(ctx context.Context, eventID CommunicationEventID, payload CommunicationPayload) error
}

func (c *communicationsDispatcher) SendCommunicationEvent(ctx context.Context, eventID CommunicationEventID, payload CommunicationPayload) error {
	commTemplates := *c.communications.getCommunicationTemplates()
	eventDestinations := commTemplates[eventID]
	if len(eventDestinations) == 0 {
		c.logger.Error(ctx, fmt.Sprintf("no_destination_registered_for_eventId=%s", eventID))
		return errors.New(string(CommErrNoDestinationFound))
	}
	for _, destination := range eventDestinations {
		c.logger.Info(ctx, fmt.Sprintf("attempting_to_send_%s_event", destination), map[string]interface{}{
			"eventId": eventID,
			"payload": payload,
		})
		switch destination {
		case CommEventTypeSSE:
			err := c.sendSSEEvent(ctx, eventID, payload)
			if err != nil {
				c.logger.Error(ctx, fmt.Sprintf("error_occurred_while_processing_sse_event... error=%s ; payload=%s", err, payload))
			}
			break
		case CommEventTypeGlobalAlert:
			break
		case CommEventTypeAudit:
			break
		default:
			c.logger.Warn(ctx, fmt.Sprintf("unknown_event_destination_encountered; destination=%s; skipping_event...", destination))
			return errors.New(string(CommErrUnknownCommDestination))
		}
		c.logger.Info(ctx, fmt.Sprintf("successfully_sent_%s_event", destination))
	}
	return nil
}

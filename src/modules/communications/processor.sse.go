/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

import (
	"context"
	"lruCache/poc/src/pkg/kafka"

	"github.com/pkg/errors"
)

func (c *communicationsDispatcher) sendSSEEvent(ctx context.Context, eventID CommunicationEventID, payload CommunicationPayload) error {
	sseEventBuilder := c.communications.getSSEEventBuilder()
	if sseEventBuilder == nil {
		errMsg := string(CommErrSSEEventBuilderMethodNotFound)
		c.logger.Error(ctx, "no_builder_method_defined_for_sse_events...terminating__sendSSEEvent__", errMsg)
		return errors.New(errMsg)
	}
	message, err := sseEventBuilder(ctx, CommUserDetails{AppUserID: c.appUserID}, eventID, payload)
	if err != nil {
		errMsg := string(CommErrSSEEventBuildFailed)
		c.logger.Error(ctx, "error_while_building_sse_event_payloads", errors.Wrap(err, errMsg))
		return errors.New(errMsg)
	}
	c.logger.Debug(ctx, "sse_payload_successfully_built...attempting_send")
	kp := kafka.GetKafkaProducerFromCache(c.communications.getAppName(), c.logger)
	err = kp.ProduceKafkaMessage(ctx, c.communications.getSSEEventDestination(), message)
	if err != nil {
		c.logger.Error(ctx, "failed_to_send_sse_event", err, map[string]interface{}{
			"appUserId":    c.appUserID,
			"eventId":      eventID,
			"payload":      payload,
			"kafkaPayload": message,
		})
		return err
	}
	c.logger.Info(ctx, "successfully_sent_sse_events", message.Events, message.AppUserID)
	return nil
}

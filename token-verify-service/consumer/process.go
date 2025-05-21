package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"tvs/dbops/model"

	"github.com/segmentio/kafka-go"
)

// TODO: imlement the database operation interface properly
// processMessage processes a Kafka message and stores the token in the database.
func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	var message model.KafkaMessage
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return fmt.Errorf("invalid message format")
	}

	reqID := message.ReqID

	if message.HashKey == "" {
		log.Printf("[request_id=%s] Hashkey is required", reqID)
		return fmt.Errorf("hashkey is required")
	}

	token := model.Token{
		Hashkey:            message.HashKey,
		ExpiresAt:          message.ExpiresAt,
		Disabled:           message.Disabled,
		RateLimitPerMinute: message.RateLimitPerMin,
	}

	// Keeping the event is seprate types to scale it later

	switch message.EventType {
	case "event.create":
		if err := c.tokenOps.Save(ctx, reqID, &token); err != nil {
			log.Printf("[request_id=%s] Failed to create token: %v", reqID, err)
			return fmt.Errorf("failed to create token")
		}
		log.Printf("[request_id=%s] Token created successfully: %v", reqID, token)
	case "event.update":
		if err := c.tokenOps.Save(ctx, reqID, &token); err != nil {
			log.Printf("[request_id=%s] Failed to update token: %v", reqID, err)
			return fmt.Errorf("failed to update token")
		}
		log.Printf("[request_id=%s] Token updated successfully: %v", reqID, token)
	case "event.disable":
		disableMessage := model.DisableMessage{
			HashKey:  message.HashKey,
			Disabled: message.Disabled,
			ReqID:    reqID,
		}
		if err := c.tokenOps.DisableToken(ctx, disableMessage); err != nil {
			log.Printf("[request_id=%s] Failed to disable token: %v", reqID, err)
			return fmt.Errorf("failed to disable token")
		}
		log.Printf("[request_id=%s] Token disable successfully: %v", reqID, token)
	}

	return nil
}

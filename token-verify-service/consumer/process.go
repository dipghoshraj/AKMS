package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"tvs/dbops"
	"tvs/dbops/model"
	"tvs/store"

	"github.com/segmentio/kafka-go"
)

// TODO: imlement the database operation interface properly
// processMessage processes a Kafka message and stores the token in the database.
func processMessage(ctx context.Context, reqID string, msg kafka.Message) error {
	var token model.Token
	if err := json.Unmarshal(msg.Value, &token); err != nil {
		log.Printf("[request_id=%s] JSON unmarshal error: %v", reqID, err)
		return fmt.Errorf("invalid message format")
	}

	if token.Hashkey == "" {
		log.Printf("[request_id=%s] Hashkey is required", reqID)
		return fmt.Errorf("hashkey is required")
	}

	dbops := dbops.NewOpsManager(store.DataBase)
	if err := dbops.TokenOps.Create(ctx, &token); err != nil {
		log.Printf("[request_id=%s] Database operation error: %v", reqID, err)
		return fmt.Errorf("failed to create token in database")
	}
	return nil
}

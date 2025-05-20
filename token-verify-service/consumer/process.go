package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"tvs/dbops/model"

	"github.com/segmentio/kafka-go"
)

func processMessage(ctx context.Context, reqID string, msg kafka.Message) error {
	var token model.Token
	if err := json.Unmarshal(msg.Value, &token); err != nil {
		log.Printf("[request_id=%s] JSON unmarshal error: %v", reqID, err)
		return fmt.Errorf("invalid message format")
	}

}

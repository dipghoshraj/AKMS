package dbops

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"tvs/dbops/model"
)

type TokenOps interface {
	Save(ctx context.Context, reqID string, input *model.Token) error
	DisableToken(ctx context.Context, disableMesage model.DisableMessage) error
}

func (t *tokenOps) Save(ctx context.Context, reqID string, token *model.Token) error {
	if err := t.db.WithContext(ctx).Save(token).Error; err != nil {
		log.Printf("[request_id=%s] Failed to update token: %v", reqID, err)
		return err
	}

	go func() {
		if err := t.SetRedis(ctx, token.Hashkey, reqID); err != nil {
			log.Printf("[request_id=%s] Failed to set token in Redis: %v", reqID, err)
		}
	}()
	return nil
}

func (t *tokenOps) DisableToken(ctx context.Context, disableMesage model.DisableMessage) error {

	result := t.db.WithContext(ctx).Model(&model.Token{}).Where("hashkey = ?", disableMesage.HashKey).Update("disabled", disableMesage.Disabled)
	if result.Error != nil {
		log.Printf("[request_id=%s] Failed to update token: %v", disableMesage.ReqID, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("[request_id=%s] Failed to update token: %v", disableMesage.ReqID, "token not found")
		return fmt.Errorf("token not found")
	}

	go func() {
		if err := t.SetRedis(ctx, disableMesage.HashKey, disableMesage.ReqID); err != nil {
			log.Printf("[request_id=%s] Failed to set token in Redis: %v", disableMesage.ReqID, err)
		}
	}()

	return nil
}

func (t *tokenOps) SetRedis(ctx context.Context, key string, reqID string) error {
	var token model.Token
	if err := t.db.WithContext(ctx).Where("hashkey = ?", key).First(&token).Error; err != nil {
		log.Printf("[request_id=%s] Failed to find token: %v", key, err)
		return err
	}

	redisValue := map[string]any{
		"disabled":   token.Disabled,
		"rate_limit": token.RateLimitPerMinute,
		"expired_at": token.ExpiresAt}

	value, err := json.Marshal(redisValue)
	if err != nil {
		log.Printf("[request_id=%s] Failed to marshal token: %v", reqID, err)
		return err
	}

	ttl := 5 * time.Minute
	if err := t.redis.Set(ctx, key, value, ttl).Err(); err != nil {
		log.Printf("[request_id=%s] Failed to set token in Redis: %v", key, err)
		return err
	}

	log.Printf("[request_id=%s] Token set in Redis: %s", reqID, key)

	return nil
}

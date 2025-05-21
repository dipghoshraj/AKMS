package dbops

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"tvs/dbops/model"
)

type TokenOps interface {
	Save(ctx context.Context, reqID string, input *model.Token) error
	DisableToken(ctx context.Context, disableMesage model.DisableMessage) error
	GetRedisToken(ctx context.Context, key string) (bool, error)
}

func (t *tokenOps) Save(ctx context.Context, reqID string, token *model.Token) error {
	if err := t.db.WithContext(ctx).Save(token).Error; err != nil {
		log.Printf("[request_id=%s] Failed to update token: %v", reqID, err)
		return err
	}

	go func() {
		if err := t.setRedis(ctx, token.Hashkey, reqID); err != nil {
			log.Printf("[request_id=%s] Failed to set token in Redis: %v", reqID, err)
		}

		rateKey := fmt.Sprintf("rate_limit:%s", token.Hashkey)
		_, err := t.redis.Incr(ctx, rateKey).Result()
		if err != nil {
			log.Printf("[request_id=%s] Failed to increment rate limit: %v", reqID, err)
			return
		}
		t.redis.Expire(ctx, rateKey, time.Minute)
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
		if err := t.setRedis(ctx, disableMesage.HashKey, disableMesage.ReqID); err != nil {
			log.Printf("[request_id=%s] Failed to set token in Redis: %v", disableMesage.ReqID, err)
		}
	}()

	return nil
}

func (t *tokenOps) setRedis(ctx context.Context, key string, reqID string) error {
	var token model.Token
	if err := t.db.WithContext(ctx).Where("hashkey = ?", key).First(&token).Error; err != nil {
		log.Printf("[request_id=%s] Failed to find token: %v", key, err)
		return err
	}

	redisValue := model.RedisMeta{
		Disabled:  token.Disabled,
		ExpiresAt: token.ExpiresAt,
		RateLimit: token.RateLimitPerMinute,
	}

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

func (t *tokenOps) GetRedisToken(ctx context.Context, key string) (bool, error) {
	hash := sha256.Sum256([]byte(key))
	hashkey := hex.EncodeToString(hash[:])
	var meta model.RedisMeta
	var token model.Token

	cached, err := t.redis.Get(ctx, hashkey).Result()
	if err == nil {

		if err := json.Unmarshal([]byte(cached), &meta); err == nil {
			if meta.Disabled || meta.ExpiresAt.Before(time.Now()) {
				log.Printf("[request_id=%s] Token is disabled or expired: %s", key, err)
				return false, fmt.Errorf("token is disable or expired")
			}
			return true, nil
		}
	}

	log.Printf("[request_id=%s] Failed to get token from Redis: %v", key, err)
	if err := t.db.WithContext(ctx).Where("hashkey = ?", key).First(&token).Error; err != nil {
		log.Printf("[request_id=%s] Failed to find token: %v", key, err)
		return false, fmt.Errorf("token not found")
	}

	if token.Disabled || token.ExpiresAt.Before(time.Now()) {
		log.Printf("[request_id=%s] Token is disabled or expired: %s", key, err)
		return false, fmt.Errorf("token is disable or expired")
	}

	rateExceed, err := t.CheckRateLimit(ctx, key, token.RateLimitPerMinute)
	if err != nil {
		log.Printf("[request_id=%s] Failed to check rate limit: %v", key, err)
		return false, err
	}

	if rateExceed {
		log.Printf("[request_id=%s] Rate limit exceeded for token: %s", key, err)
		return false, fmt.Errorf("rate limit exceeded")
	}

	return true, nil
}

func (t *tokenOps) CheckRateLimit(ctx context.Context, key string, limit int64) (bool, error) {
	rateKey := fmt.Sprintf("rate_limit:%s", key)
	rateLimit, err := t.redis.Incr(ctx, rateKey).Result()

	if err != nil {
		log.Printf("[request_id=%s] Failed to get rate limit: %v", key, err)
		return false, err
	}

	if rateLimit > limit {
		return true, nil
	}

	return false, nil
}

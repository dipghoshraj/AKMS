package dbops

import (
	"akm/dbops/model"
	"akm/dbops/producer"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenOps interface {
	Create(ctx context.Context, input *model.TokenCreateInput) (*model.Token, error)
	GetAll(ctx context.Context) ([]*model.Token, error)
	GetByKey(ctx context.Context, key string) (*model.Token, error)
	Disable(ctx context.Context, key string) error
	Update(ctx context.Context, key string, input *model.TokenUpdateInput) (*model.Token, error)
	Delete(ctx context.Context, key string) error
}

func generateHash(tokenKey string) string {
	hash := sha256.Sum256([]byte(tokenKey))
	return hex.EncodeToString(hash[:])
}

func (to *tokenOps) Create(ctx context.Context, input *model.TokenCreateInput) (*model.Token, error) {

	hashkey, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err)
		return nil, err
	}

	hash := generateHash(hashkey.String())
	expiresAt := time.Now().Add(time.Duration(input.ExpiresAt) * time.Minute)

	token := &model.Token{
		Hashkey:            hash,
		RateLimitPerMinute: input.RateLimitPerMinute,
		ExpiresAt:          expiresAt,
		Disabled:           false,
	}

	err = to.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// DB insert
		if err := tx.Create(token).Error; err != nil {
			return fmt.Errorf("error creating token: %w", err)
		}

		reqID, ok := ctx.Value("reqID").(string)
		if !ok {
			fmt.Println("Request ID not found in context")
			return fmt.Errorf("request ID not found in context")
		}

		// message := map[string]string{
		// 	"hashkey":            hash,
		// 	"rate_limit_per_min": fmt.Sprintf("%d", input.RateLimitPerMinute),
		// 	"expires_at":         expiresAt.Format(time.RFC3339),
		// 	"disabled":           fmt.Sprintf("%t", token.Disabled),
		// }

		message := model.KafkaMessage{
			HashKey:         hash,
			RateLimitPerMin: token.RateLimitPerMinute,
			ExpiresAt:       expiresAt,
			Disabled:        token.Disabled,
		}

		if err = producer.NewProducer().PushMessage(reqID, message); err != nil {
			fmt.Printf("Error pushing message to Kafka: %v\n", err)
			return fmt.Errorf("error pushing message to Kafka: %w", err)
		}

		return nil // commit the transaction
	})

	if err != nil {
		return nil, err // transaction rolled back
	}
	token.Hashkey = hashkey.String()

	return token, nil
}

func (to *tokenOps) GetAll(ctx context.Context) ([]*model.Token, error) {
	var tokens []*model.Token
	if err := to.db.WithContext(ctx).Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

func (to *tokenOps) GetByKey(ctx context.Context, key string) (*model.Token, error) {
	var token model.Token

	hash := generateHash(key)
	if err := to.db.WithContext(ctx).Where("hashkey = ?", hash).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (to *tokenOps) Disable(ctx context.Context, key string) error {
	hash := generateHash(key)

	result := to.db.WithContext(ctx).Model(&model.Token{}).Where("hashkey = ?", hash).Update("disabled", true)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

func (to *tokenOps) Update(ctx context.Context, key string, input *model.TokenUpdateInput) (*model.Token, error) {
	hash := generateHash(key)

	// Perform the update
	result := to.db.WithContext(ctx).Model(&model.Token{}).Where("hashkey = ?", hash).Updates(input)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("token not found")
	}

	// Fetch and return the updated token
	var token model.Token
	if err := to.db.WithContext(ctx).Where("hashkey = ?", hash).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (to *tokenOps) Delete(ctx context.Context, key string) error {
	hash := generateHash(key)

	result := to.db.WithContext(ctx).Where("hashkey = ?", hash).Delete(&model.Token{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

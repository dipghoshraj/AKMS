package dbops

import (
	"akm/dbops/model"
	"akm/store"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	if err := to.db.WithContext(ctx).Create(token).Error; err != nil {
		return nil, err
	}

	if err := store.DataBase.Find(&token).Error; err != nil {
		return nil, err
	}

	token.Hashkey = hashkey.String() // Overwrite just for return

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

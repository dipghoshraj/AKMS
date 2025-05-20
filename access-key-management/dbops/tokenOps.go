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

	// if err := store.DataBase.WithContext(ctx).Create(token).Error; err != nil {
	// 	return nil, err
	// }

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
	if err := to.db.WithContext(ctx).Where("hashkey = ?", key).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

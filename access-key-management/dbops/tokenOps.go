package dbops

import (
	"akm/dbops/model"
	"akm/store"
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type TokenOps interface {
	Create(ctx context.Context, input *model.TokenCreateInput) (*model.Token, error)
}

func generateHash(tokenKey string) string {
	hash := sha256.Sum256([]byte(tokenKey))
	return hex.EncodeToString(hash[:])
}

func (to *tokenOps) Create(ctx context.Context, input *model.TokenCreateInput) (*model.Token, error) {

	hash := generateHash(input.Hashkey)

	token := &model.Token{
		Hashkey:            hash,
		RateLimitPerMinute: input.RateLimitPerMinute,
		ExpiresAt:          input.ExpiresAt,
		Disabled:           false,
	}

	if err := store.DataBase.WithContext(ctx).Create(token).Error; err != nil {
		return nil, err
	}

	if err := store.DataBase.Find(&token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

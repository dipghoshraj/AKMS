package dbops

import (
	"context"
	"tvs/dbops/model"
)

type TokenOps interface {
	Create(ctx context.Context, input *model.Token) error
}

func (t *tokenOps) Create(ctx context.Context, input *model.Token) error {
	token := &model.Token{
		Hashkey:            input.Hashkey,
		RateLimitPerMinute: input.RateLimitPerMinute,
		ExpiresAt:          input.ExpiresAt,
		Disabled:           input.Disabled,
	}

	if err := t.db.WithContext(ctx).Create(token).Error; err != nil {
		return err
	}

	return nil
}

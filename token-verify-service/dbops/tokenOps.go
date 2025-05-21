package dbops

import (
	"context"
	"fmt"
	"log"
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

	return nil
}

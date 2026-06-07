package repository

import (
	"context"

	"gorm.io/gorm"

	"secureops/backend-go/api/model"
)

type WafEventRepository struct {
	db *gorm.DB
}

func NewWafEventRepository(db *gorm.DB) *WafEventRepository {
	return &WafEventRepository{db: db}
}

func (r *WafEventRepository) Save(ctx context.Context, event model.WafEvent) error {
	return r.db.WithContext(ctx).Create(&event).Error
}

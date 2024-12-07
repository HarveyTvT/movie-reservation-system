package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type OrderRepository struct {
	db *bun.DB
}

func NewOrderRepository(db *bun.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, record *model.Order) error {
	_, err := r.db.NewInsert().Model(record).Exec(ctx)
	return err
}

func (r *OrderRepository) Delete(ctx context.Context, record *model.Order) error {
	_, err := r.db.NewDelete().Model(record).WherePK().Exec(ctx)
	return err
}

func (r *OrderRepository) Update(ctx context.Context, record *model.Order, columns ...string) error {
	_, err := r.db.NewUpdate().Model(record).WherePK().Column(columns...).Exec(ctx, &record)
	return err
}

func (r *OrderRepository) ListByUserId(ctx context.Context, userId uint64, offset uint64, limit uint64) ([]*model.Order, uint64, error) {
	results := make([]*model.Order, 0)
	cnt, err := r.db.NewSelect().Model(&model.Order{}).
		Where("user_id = ?", userId).
		Offset(int(offset)).
		Limit(int(limit)).
		Order("created_at DESC").
		ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, err
	}

	return results, uint64(cnt), nil
}

package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type HallRepository struct {
	db *bun.DB
}

func NewHallRepository(db *bun.DB) *HallRepository {
	return &HallRepository{db}
}

func (r *HallRepository) Create(ctx context.Context, hall *model.Hall) error {
	_, err := r.db.NewInsert().Model(hall).Exec(ctx)
	return err
}

func (r *HallRepository) ListByTheaterId(ctx context.Context, theaterId uint64, offset uint64, limit uint64) ([]*model.Hall, uint64, error) {
	results := make([]*model.Hall, 0)
	cnt, err := r.db.NewSelect().Model(&model.Hall{}).
		Where("theater_id = ?", theaterId).
		Offset(int(offset)).Limit(int(limit)).ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, err
	}

	return results, uint64(cnt), nil
}

func (r *HallRepository) Update(ctx context.Context, hall *model.Hall, columns ...string) error {
	_, err := r.db.NewUpdate().Model(hall).WherePK().Column(columns...).Exec(ctx)
	return err
}

func (r *HallRepository) Delete(ctx context.Context, hall *model.Hall) error {
	_, err := r.db.NewDelete().Model(hall).WherePK().Exec(ctx)
	return err
}

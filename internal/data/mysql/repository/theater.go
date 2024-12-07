package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type TheaterRepository struct {
	db *bun.DB
}

func NewTheaterRepository(db *bun.DB) *TheaterRepository {
	return &TheaterRepository{db}
}

func (r *TheaterRepository) Create(ctx context.Context, theater *model.Theater) error {
	_, err := r.db.NewInsert().Model(theater).Exec(ctx)
	return err
}

func (r *TheaterRepository) Get(ctx context.Context, id uint64) (*model.Theater, error) {
	theater := new(model.Theater)
	err := r.db.NewSelect().Model(theater).Where("id = ?", id).Scan(ctx, &theater)
	if err != nil {
		return nil, err
	}
	return theater, nil
}

func (r *TheaterRepository) List(ctx context.Context, offset uint64, limit uint64) ([]*model.Theater, uint64, error) {
	results := make([]*model.Theater, 0)
	cnt, err := r.db.NewSelect().Model(&model.Theater{}).Offset(int(offset)).Limit(int(limit)).ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, err
	}

	return results, uint64(cnt), nil
}

func (r *TheaterRepository) Update(ctx context.Context, theater *model.Theater, columns ...string) error {
	_, err := r.db.NewUpdate().Model(theater).WherePK().Column(columns...).Exec(ctx)
	return err
}

func (r *TheaterRepository) Delete(ctx context.Context, theater *model.Theater) error {
	_, err := r.db.NewDelete().Model(theater).WherePK().Exec(ctx)
	return err
}


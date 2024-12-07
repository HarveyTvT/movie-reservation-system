package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type ShowRepository struct {
	db *bun.DB
}

func NewShowRepository(db *bun.DB) *ShowRepository {
	return &ShowRepository{db}
}

func (r *ShowRepository) Create(ctx context.Context, show *model.Show) error {
	_, err := r.db.NewInsert().Model(show).Exec(ctx)
	return err
}

func (r *ShowRepository) Get(ctx context.Context, id uint64) (*model.Show, error) {
	show := new(model.Show)
	err := r.db.NewSelect().Model(show).Where("id = ?", id).Scan(ctx, &show)
	if err != nil {
		return nil, err
	}
	return show, nil
}

func (r *ShowRepository) ListByMovie(ctx context.Context, movieID uint64, offset uint64, limit uint64) ([]*model.Show, uint64, error) {
	results := make([]*model.Show, 0)
	cnt, err := r.db.NewSelect().Model(&model.Show{}).Where("movie_id = ?", movieID).Offset(int(offset)).Limit(int(limit)).ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, err
	}

	return results, uint64(cnt), nil
}

func (r *ShowRepository) Update(ctx context.Context, show *model.Show, columns ...string) error {
	_, err := r.db.NewUpdate().Model(show).WherePK().Column(columns...).Exec(ctx)
	return err
}

func (r *ShowRepository) Delete(ctx context.Context, show *model.Show) error {
	_, err := r.db.NewDelete().Model(show).WherePK().Exec(ctx)
	return err
}

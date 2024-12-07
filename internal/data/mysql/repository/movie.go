package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type MovieRepository struct {
	db *bun.DB
}

func NewMovieRepository(db *bun.DB) *MovieRepository {
	return &MovieRepository{db}
}

func (r *MovieRepository) Create(ctx context.Context, movie *model.Movie) error {
	_, err := r.db.NewInsert().Model(movie).Exec(ctx)
	return err
}

func (r *MovieRepository) Get(ctx context.Context, id uint64) (*model.Movie, error) {
	movie := new(model.Movie)
	err := r.db.NewSelect().Model(movie).Where("id = ?", id).Scan(ctx, &movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (r *MovieRepository) List(ctx context.Context, offset uint64, limit uint64) ([]*model.Movie, uint64, error) {
	results := make([]*model.Movie, 0)
	cnt, err := r.db.NewSelect().Model(&model.Movie{}).Offset(int(offset)).Limit(int(limit)).ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, err
	}

	return results, uint64(cnt), nil
}

func (r *MovieRepository) Update(ctx context.Context, movie *model.Movie, columns ...string) error {
	_, err := r.db.NewUpdate().Model(movie).WherePK().Column(columns...).Exec(ctx)
	return err
}

func (r *MovieRepository) MGet(ctx context.Context, ids []uint64) ([]*model.Movie, error) {
	results := make([]*model.Movie, 0)
	if len(ids) == 0 {
		return results, nil
	}
	err := r.db.NewSelect().Model(&results).Where("id IN (?)", bun.In(ids)).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return results, nil
}

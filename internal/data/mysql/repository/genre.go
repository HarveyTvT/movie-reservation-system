package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type GenreRepository struct {
	db *bun.DB
}

func NewGenreRepository(db *bun.DB) *GenreRepository {
	return &GenreRepository{db}
}

func (r *GenreRepository) Create(ctx context.Context, genre *model.Genre) error {
	_, err := r.db.NewInsert().Model(genre).Ignore().Exec(ctx)
	return err
}

func (r *GenreRepository) List(ctx context.Context, offset uint64, limit uint64) ([]*model.Genre, uint64, error) {
	results := make([]*model.Genre, 0)
	cnt, err := r.db.NewSelect().Model(&model.Genre{}).Offset(int(offset)).Limit(int(limit)).ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, err
	}

	return results, uint64(cnt), nil
}

func (r *GenreRepository) Delete(ctx context.Context, genre string) error {
	_, err := r.db.NewDelete().Model(&model.Genre{}).Where("genre = ?", genre).Exec(ctx)
	return err
}

package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type MovieGenreRepository struct {
	db *bun.DB
}

func NewMovieGenreRepository(db *bun.DB) *MovieGenreRepository {
	return &MovieGenreRepository{db}
}

func (r *MovieGenreRepository) Create(ctx context.Context, movieGenre *model.MovieGenre) error {
	_, err := r.db.NewInsert().Model(movieGenre).Exec(ctx)
	return err
}

func (r *MovieGenreRepository) ListMovieIds(ctx context.Context, genre string, offset uint64, limit uint64) ([]uint64, uint64, error) {
	results := make([]uint64, 0)
	cnt, err := r.db.NewSelect().Column("movie_id").Model(&model.MovieGenre{}).
		Where("genre = ?", genre).
		Offset(int(offset)).Limit(int(limit)).ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, err
	}

	return results, uint64(cnt), nil
}

func (r *MovieGenreRepository) Delete(ctx context.Context, movieGenre *model.MovieGenre) error {
	_, err := r.db.NewDelete().Model(movieGenre).WherePK().Exec(ctx)
	return err
}

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
	_, err := r.db.NewInsert().Ignore().Model(movieGenre).Exec(ctx)
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
	_, err := r.db.NewDelete().Model(movieGenre).Where("movie_id = ? and genre = ?", movieGenre.MovieID, movieGenre.Genre).Exec(ctx)
	return err
}

func (r *MovieGenreRepository) MGetMovieGenres(ctx context.Context, movieIDs []uint64) (map[uint64][]string, error) {
	results := make([]model.MovieGenre, 0)
	if len(movieIDs) == 0 {
		return nil, nil
	}
	err := r.db.NewSelect().Model(&results).Where("movie_id IN (?)", bun.In(movieIDs)).Scan(ctx)
	if err != nil {
		return nil, err
	}

	movieGenres := make(map[uint64][]string)
	for _, result := range results {
		movieGenres[result.MovieID] = append(movieGenres[result.MovieID], result.Genre)
	}

	return movieGenres, nil
}

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
	result, err := r.db.NewInsert().Model(movie).Exec(ctx)
	if err != nil {
		return err
	}
	resultID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	movie.ID = uint64(resultID)
	return nil
}

func (r *MovieRepository) Get(ctx context.Context, id uint64) (*model.Movie, error) {
	movie := new(model.Movie)
	err := r.db.NewSelect().Model(movie).Where("id = ?", id).Scan(ctx, &movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (r *MovieRepository) List(ctx context.Context, genre string, offset uint64, limit uint64) ([]*model.Movie, uint64, error) {
	results := make([]*model.Movie, 0)
	s := r.db.NewSelect().Model(&model.Movie{}).Offset(int(offset)).Limit(int(limit))

	if genre != "" {
		s.Join("INNER JOIN movie_genres ON movie_genres.movie_id = movies.id").Where("movie_genres.genre = ?", genre)
	}

	cnt, err := s.ScanAndCount(ctx, &results)
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

func (r *MovieRepository) Delete(ctx context.Context, movieID uint64) error {
	_, err := r.db.NewDelete().Model(&model.Movie{ID: movieID}).WherePK().Exec(ctx)
	return err
}

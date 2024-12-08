package genre

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/repository"
	"github.com/spf13/cast"
)

type Biz interface {
	Create(ctx context.Context, genre string) error
	List(ctx context.Context, offset uint64, limit uint64) ([]string, uint64, error)
	Add(ctx context.Context, movieID string, genre string) error
	Remove(ctx context.Context, movieID string, genre string) error
}

type biz struct {
	genreRepo      *repository.GenreRepository
	movieGenreRepo *repository.MovieGenreRepository
}

func NewBiz(
	genreRepo *repository.GenreRepository,
	movieGenreRepo *repository.MovieGenreRepository,
) Biz {
	return &biz{
		genreRepo:      genreRepo,
		movieGenreRepo: movieGenreRepo,
	}
}

func (b *biz) Create(ctx context.Context, genre string) error {
	return b.genreRepo.Create(ctx, &model.Genre{Name: genre})
}

func (b *biz) List(ctx context.Context, offset uint64, limit uint64) ([]string, uint64, error) {
	genres, cnt, err := b.genreRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	results := make([]string, 0, len(genres))
	for _, genre := range genres {
		results = append(results, genre.Name)
	}

	return results, cnt, nil
}

func (b *biz) Add(ctx context.Context, movieID string, genre string) error {
	return b.movieGenreRepo.Create(ctx, &model.MovieGenre{MovieID: cast.ToUint64(movieID), Genre: genre})
}

func (b *biz) Remove(ctx context.Context, movieID string, genre string) error {
	return b.movieGenreRepo.Delete(ctx, &model.MovieGenre{MovieID: cast.ToUint64(movieID), Genre: genre})
}

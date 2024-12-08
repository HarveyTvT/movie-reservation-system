package movie

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/repository"
	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Biz interface {
	Create(ctx context.Context, request *movie_reservation.CreateMovieRequest) error
	Update(ctx context.Context, request *movie_reservation.UpdateMovieRequest) error
	Delete(ctx context.Context, movieID uint64) error
	List(ctx context.Context, request *movie_reservation.ListMoviesRequest) ([]*movie_reservation.Movie, uint64, error)
}

type biz struct {
	movieRepo *repository.MovieRepository
}

func NewBiz(movieRepo *repository.MovieRepository) Biz {
	return &biz{
		movieRepo: movieRepo,
	}
}

func (b *biz) Create(ctx context.Context, request *movie_reservation.CreateMovieRequest) error {
	record := &model.Movie{
		Name:        request.Title,
		Description: request.Description,
		Poster:      request.Poster,
		Duration:    request.Duration,
	}

	err := b.movieRepo.Create(ctx, record)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (b *biz) Update(ctx context.Context, request *movie_reservation.UpdateMovieRequest) error {
	record, err := b.movieRepo.Get(ctx, cast.ToUint64(request.Id))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	columns := make([]string, 0)

	if request.Title != nil {
		record.Name = request.Title.Value
		columns = append(columns, "name")
	}

	if request.Description != nil {
		record.Description = request.Description.Value
		columns = append(columns, "description")
	}

	if request.Poster != nil {
		record.Poster = request.Poster.Value
		columns = append(columns, "poster")
	}

	if request.Duration != nil {
		record.Duration = request.Duration.Value
		columns = append(columns, "duration")
	}

	err = b.movieRepo.Update(ctx, record, columns...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (b *biz) Delete(ctx context.Context, movieID uint64) error {
	err := b.movieRepo.Delete(ctx, movieID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (b *biz) List(ctx context.Context, request *movie_reservation.ListMoviesRequest) ([]*movie_reservation.Movie, uint64, error) {
	offset := cast.ToUint64(request.Offset)
	limit := cast.ToUint64(request.Limit)

	records, count, err := b.movieRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	movies := make([]*movie_reservation.Movie, 0, len(records))
	for _, record := range records {
		movies = append(movies, record.ToPB())
	}

	return movies, count, nil
}

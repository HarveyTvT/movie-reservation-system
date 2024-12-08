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
	Create(ctx context.Context, request *movie_reservation.CreateMovieRequest) (uint64, error)
	Update(ctx context.Context, request *movie_reservation.UpdateMovieRequest) error
	Delete(ctx context.Context, movieID string) error
	List(ctx context.Context, request *movie_reservation.ListMoviesRequest) ([]*movie_reservation.Movie, uint64, error)
}

type biz struct {
	movieRepo      *repository.MovieRepository
	genreRepo      *repository.GenreRepository
	movieGenreRepo *repository.MovieGenreRepository
}

func NewBiz(
	movieRepo *repository.MovieRepository,
	genreRepo *repository.GenreRepository,
	movieGenre *repository.MovieGenreRepository,
) Biz {
	return &biz{
		movieRepo:      movieRepo,
		genreRepo:      genreRepo,
		movieGenreRepo: movieGenre,
	}
}

func (b *biz) Create(ctx context.Context, request *movie_reservation.CreateMovieRequest) (uint64, error) {
	record := &model.Movie{
		Name:        request.Title,
		Description: request.Description,
		Poster:      request.Poster,
		Duration:    request.Duration,
	}

	err := b.movieRepo.Create(ctx, record)
	if err != nil {
		return 0, status.Error(codes.Internal, err.Error())
	}

	for _, genre := range request.Genres {
		b.genreRepo.Create(ctx, &model.Genre{Name: genre})
		b.movieGenreRepo.Create(ctx, &model.MovieGenre{
			MovieID: record.ID,
			Genre:   genre,
		})
	}

	return record.ID, nil
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

func (b *biz) Delete(ctx context.Context, movieID string) error {
	err := b.movieRepo.Delete(ctx, cast.ToUint64(movieID))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (b *biz) List(ctx context.Context, request *movie_reservation.ListMoviesRequest) ([]*movie_reservation.Movie, uint64, error) {
	offset := cast.ToUint64(request.Offset)
	limit := cast.ToUint64(request.Limit)

	records, count, err := b.movieRepo.List(ctx, request.GetGenre(), offset, limit)
	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	movieIDs := make([]uint64, 0, len(records))
	for _, record := range records {
		movieIDs = append(movieIDs, record.ID)
	}

	movieGenres, err := b.movieGenreRepo.MGetMovieGenres(ctx, movieIDs)
	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	movies := make([]*movie_reservation.Movie, 0, len(records))
	for _, record := range records {
		result := record.ToPB()
		if genres, ok := movieGenres[record.ID]; ok {
			result.Genres = genres
		}
		movies = append(movies, result)

	}

	return movies, count, nil
}

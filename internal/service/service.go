package service

import (
	"context"
	"net/http"

	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/harveytvt/movie-reservation-system/internal/auth"
	"github.com/harveytvt/movie-reservation-system/internal/biz/genre"
	"github.com/harveytvt/movie-reservation-system/internal/biz/movie"
	"github.com/harveytvt/movie-reservation-system/internal/biz/user"
	"github.com/harveytvt/movie-reservation-system/internal/config"
	"github.com/spf13/cast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type service struct {
	movie_reservation.UnimplementedMovieReservationServiceServer
	grpc_health_v1.UnimplementedHealthServer

	user  user.Biz
	movie movie.Biz
	genre genre.Biz
}

type Service interface {
	movie_reservation.MovieReservationServiceServer
	grpc_health_v1.HealthServer
}

func NewService() Service {
	return &service{
		user:  userBiz,
		movie: movieBiz,
		genre: genreBiz,
	}
}

func (s *service) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *service) Register(ctx context.Context, req *movie_reservation.RegisterRequest) (*movie_reservation.RegisterResponse, error) {
	_, err := s.user.Create(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.RegisterResponse{}, nil
}

func (s *service) Login(ctx context.Context, req *movie_reservation.LoginRequest) (*movie_reservation.LoginResponse, error) {
	u, err := s.user.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	jwtPayload := auth.JwtPayload{
		Username: u.Username,
		Role:     u.Role,
	}
	jwt := auth.NewJwt(jwtPayload, config.Get().Secret)

	cookie := http.Cookie{
		Name:   "authorization",
		Value:  jwt.String(),
		Path:   "/",
		MaxAge: 60 * 60 * 24,
	}

	grpc.SendHeader(ctx, metadata.New(map[string]string{
		"set-cookie": cookie.String(),
	}))

	return &movie_reservation.LoginResponse{
		Token: jwt.String(),
	}, nil
}

func (s *service) Whoami(ctx context.Context, req *movie_reservation.WhoamiRequest) (*movie_reservation.WhoamiResponse, error) {
	username := auth.UsernameFromContext(ctx)

	u, err := s.user.Find(ctx, username)
	if err != nil {
		return nil, err
	}

	return &movie_reservation.WhoamiResponse{
		User: u,
	}, nil
}

func (s *service) CreateMovie(ctx context.Context, req *movie_reservation.CreateMovieRequest) (*movie_reservation.CreateMovieResponse, error) {
	if err := auth.AssertRole(ctx, movie_reservation.User_ROLE_ADMIN); err != nil {
		return nil, err
	}

	movieID, err := s.movie.Create(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.CreateMovieResponse{
		Id: cast.ToString(movieID),
	}, nil
}

func (s *service) UpdateMovie(ctx context.Context, req *movie_reservation.UpdateMovieRequest) (*movie_reservation.UpdateMovieResponse, error) {
	if err := auth.AssertRole(ctx, movie_reservation.User_ROLE_ADMIN); err != nil {
		return nil, err
	}

	err := s.movie.Update(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.UpdateMovieResponse{}, nil
}

func (s *service) DeleteMovie(ctx context.Context, req *movie_reservation.DeleteMovieRequest) (*movie_reservation.DeleteMovieResponse, error) {
	if err := auth.AssertRole(ctx, movie_reservation.User_ROLE_ADMIN); err != nil {
		return nil, err
	}

	err := s.movie.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.DeleteMovieResponse{}, nil
}

func (s *service) ListGenres(ctx context.Context, req *movie_reservation.ListGenresRequest) (*movie_reservation.ListGenresResponse, error) {
	genres, cnt, err := s.genre.List(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.ListGenresResponse{
		Genres: genres,
		Total:  cnt,
	}, nil
}

func (s *service) CreateGenre(ctx context.Context, req *movie_reservation.CreateGenreRequest) (*movie_reservation.CreateGenreResponse, error) {
	if err := auth.AssertRole(ctx, movie_reservation.User_ROLE_ADMIN); err != nil {
		return nil, err
	}

	err := s.genre.Create(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.CreateGenreResponse{}, nil
}

func (s *service) AddMovieGenre(ctx context.Context, req *movie_reservation.AddMovieGenreRequest) (*movie_reservation.AddMovieGenreResponse, error) {
	if err := auth.AssertRole(ctx, movie_reservation.User_ROLE_ADMIN); err != nil {
		return nil, err
	}

	err := s.genre.Add(ctx, req.Id, req.Genre)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.AddMovieGenreResponse{}, nil
}

func (s *service) RemoveMovieGenre(ctx context.Context, req *movie_reservation.RemoveMovieGenreRequest) (*movie_reservation.RemoveMovieGenreResponse, error) {
	if err := auth.AssertRole(ctx, movie_reservation.User_ROLE_ADMIN); err != nil {
		return nil, err
	}

	err := s.genre.Remove(ctx, req.Id, req.Genre)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.RemoveMovieGenreResponse{}, nil
}

func (s *service) ListMovies(ctx context.Context, req *movie_reservation.ListMoviesRequest) (*movie_reservation.ListMoviesResponse, error) {
	movies, cnt, err := s.movie.List(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &movie_reservation.ListMoviesResponse{
		Movies: movies,
		Total:  cnt,
	}, nil
}

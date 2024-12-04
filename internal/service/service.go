package service

import (
	"context"
	"net/http"

	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/harveytvt/movie-reservation-system/internal/auth"
	"github.com/harveytvt/movie-reservation-system/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type service struct {
	movie_reservation.UnimplementedMovieReservationServiceServer
}

type Service interface {
	movie_reservation.MovieReservationServiceServer
}

func NewService() Service {
	return &service{}
}

func (s *service) Login(ctx context.Context, req *movie_reservation.LoginRequest) (*movie_reservation.LoginResponse, error) {
	jwt := auth.NewJwt(req.Username, config.Get().Secret)

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
	return &movie_reservation.WhoamiResponse{
		Username: username,
	}, nil
}

package service

import (
	"context"
	"net/http"

	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/harveytvt/movie-reservation-system/internal/auth"
	"github.com/harveytvt/movie-reservation-system/internal/biz/user"
	"github.com/harveytvt/movie-reservation-system/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type service struct {
	movie_reservation.UnimplementedMovieReservationServiceServer
	grpc_health_v1.UnimplementedHealthServer

	user user.Biz
}

type Service interface {
	movie_reservation.MovieReservationServiceServer
	grpc_health_v1.HealthServer
}

func NewService() Service {
	return &service{
		user: userBiz,
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

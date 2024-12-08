package user

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/harveytvt/movie-reservation-system/internal/config"
	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Biz interface {
	Find(ctx context.Context, username string) (*movie_reservation.User, error)
	Get(ctx context.Context, id uint64) (*movie_reservation.User, error)
	Create(ctx context.Context, username string, password string) (*movie_reservation.User, error)
	Login(ctx context.Context, username string, password string) (*movie_reservation.User, error)
}

type biz struct {
	user *repository.UserRepository
}

func NewBiz(
	user *repository.UserRepository,
) Biz {
	return &biz{
		user: user,
	}
}

func (b *biz) Find(ctx context.Context, username string) (*movie_reservation.User, error) {
	u, err := b.user.Find(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, err
	}

	return u.ToPB(), nil
}

func (b *biz) Get(ctx context.Context, id uint64) (*movie_reservation.User, error) {
	u, err := b.user.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, err
	}

	return u.ToPB(), nil
}

func (b *biz) Create(ctx context.Context, username string, password string) (*movie_reservation.User, error) {
	exists, err := b.user.Exists(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}

	u := &model.User{
		Username: username,
		Password: hashPassword(password),
		Role:     uint32(movie_reservation.User_ROLE_USER),
	}

	err = b.user.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return u.ToPB(), nil
}

func (b *biz) Login(ctx context.Context, username string, password string) (*movie_reservation.User, error) {
	u, err := b.user.Find(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, err
	}

	if u.Password != hashPassword(password) {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	return u.ToPB(), nil
}

func hashPassword(password string) string {
	if password == "" {
		return ""
	}

	hasher := sha256.New()
	hasher.Write([]byte(password + config.Get().Salt))

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

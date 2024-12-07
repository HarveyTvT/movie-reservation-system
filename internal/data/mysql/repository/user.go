package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Find(ctx context.Context, unsername string) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().Model(user).Where("username = ?", unsername).Scan(ctx, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *model.User, columns ...string) error {
	_, err := r.db.NewUpdate().Model(user).WherePK().Column(columns...).Exec(ctx)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, user *model.User) error {
	_, err := r.db.NewDelete().Model(user).WherePK().Exec(ctx)
	return err
}

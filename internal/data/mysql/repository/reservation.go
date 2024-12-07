package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type ReservationRepository struct {
	db *bun.DB
}

func NewReservationRepository(db *bun.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(ctx context.Context, reservation *model.Reservation) error {
	_, err := r.db.NewInsert().Model(reservation).Exec(ctx)
	return err
}

func (r *ReservationRepository) Update(ctx context.Context, reservation *model.Reservation, columns ...string) error {
	_, err := r.db.NewUpdate().Model(reservation).WherePK().Column(columns...).Exec(ctx)
	return err
}

func (r *ReservationRepository) Delete(ctx context.Context, reservation *model.Reservation) error {
	_, err := r.db.NewDelete().Model(reservation).WherePK().Exec(ctx)
	return err
}

func (r *ReservationRepository) GetByShowID(ctx context.Context, showID uint64) ([]*model.Reservation, error) {
	reservations := make([]*model.Reservation, 0)
	err := r.db.NewSelect().Model(&reservations).Where("show_id = ?", showID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

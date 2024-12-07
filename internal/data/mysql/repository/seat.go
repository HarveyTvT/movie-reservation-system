package repository

import (
	"context"

	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/model"
	"github.com/uptrace/bun"
)

type SeatRepository struct {
	db *bun.DB
}

func NewSeatRepository(db *bun.DB) *SeatRepository {
	return &SeatRepository{
		db: db,
	}
}

func (r *SeatRepository) Create(ctx context.Context, seat *model.Seat) error {
	_, err := r.db.NewInsert().Model(seat).Exec(ctx)
	return err
}

func (r *SeatRepository) Update(ctx context.Context, seat *model.Seat, columns ...string) error {
	_, err := r.db.NewUpdate().Model(seat).WherePK().Column(columns...).Exec(ctx)
	return err
}

func (r *SeatRepository) Delete(ctx context.Context, seat *model.Seat) error {
	_, err := r.db.NewDelete().Model(seat).WherePK().Exec(ctx)
	return err
}

func (r *SeatRepository) GetByHallID(ctx context.Context, hallID uint64) ([]*model.Seat, error) {
	seats := make([]*model.Seat, 0)
	err := r.db.NewSelect().Model(&seats).Where("hall_id = ?", hallID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

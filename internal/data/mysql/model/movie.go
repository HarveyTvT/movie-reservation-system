package model

import (
	"time"

	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/spf13/cast"
	"github.com/uptrace/bun"
)

type Movie struct {
	bun.BaseModel `bun:"movies,alias:movies"`

	ID          uint64     `bun:",pk"`
	CreatedAt   *time.Time `bun:"created_at"`
	UpdatedAt   *time.Time `bun:"updated_at"`
	Name        string     `bun:"name"`
	Description string     `bun:"description"`
	Poster      string     `bun:"poster"`
	Duration    uint64     `bun:"duration"`
}

func (m *Movie) ToPB() *movie_reservation.Movie {
	if m == nil {
		return nil
	}
	return &movie_reservation.Movie{
		Id:          cast.ToString(m.ID),
		Title:       m.Name,
		Description: m.Description,
		Poster:      m.Poster,
		Duration:    m.Duration,
	}
}

package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Show struct {
	bun.BaseModel `bun:"shows,alias:shows"`

	ID        uint64     `bun:",pk"`
	MovieID   uint64     `bun:"movie_id"`
	HallID    uint64     `bun:"hall_id"`
	StartedAt *time.Time `bun:"started_at"`
}

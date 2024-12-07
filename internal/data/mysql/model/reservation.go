package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Reservation struct {
	bun.BaseModel `bun:"reservations,alias:reservations"`

	ID        uint64     `bun:",pk"`
	CreatedAt *time.Time `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
	ShowID    uint64     `bun:"show_id"`
	SeatID    uint64     `bun:"seat_id"`
	Status    uint32     `bun:"status"` // 0-unknown, 1-reserved, 2-free
}

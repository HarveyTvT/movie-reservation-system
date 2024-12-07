package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"orders,alias:orders"`

	ID        uint64     `bun:",pk"`
	CreatedAt *time.Time `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
	UserID    uint64     `bun:"user_id"`
	ShowID    uint64     `bun:"show_id"`
	TheaterID uint64     `bun:"theater_id"`
	HallID    uint64     `bun:"hall_id"`
	SeatID    uint64     `bun:"seat_id"`
	Status    uint32     `bun:"status"` // 0-unknown, 1-pending, 2-confirmed, 3-canceled, 4-consumed, 5-expired
}

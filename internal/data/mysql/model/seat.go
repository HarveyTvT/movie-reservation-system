package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Seat struct {
	bun.BaseModel `bun:"seats,alias:seats"`

	ID        uint64     `bun:",pk"`
	CreatedAt *time.Time `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
	HallID    uint64     `bun:"hall_id"`
	Row       uint32     `bun:"row"`
	Col       uint32     `bun:"col"`
}

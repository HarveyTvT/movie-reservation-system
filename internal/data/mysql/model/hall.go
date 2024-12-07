package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Hall struct {
	bun.BaseModel `bun:"halls,alias:halls"`

	ID        uint64     `bun:",pk"`
	CreatedAt *time.Time `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
	TheaterID uint64     `bun:"theater_id"`
	Name      string     `bun:"name"`
}

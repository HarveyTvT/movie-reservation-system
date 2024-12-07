package model

import (
	"time"

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

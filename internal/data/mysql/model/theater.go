package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Theater struct {
	bun.BaseModel `bun:"theaters,alias:theaters"`

	ID          uint64     `bun:",pk"`
	CreatedAt   *time.Time `bun:"created_at"`
	UpdatedAt   *time.Time `bun:"updated_at"`
	Name        string     `bun:"name"`
	Description string     `bun:"description"`
}

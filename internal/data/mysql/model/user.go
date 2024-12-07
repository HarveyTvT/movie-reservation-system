package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"users,alias:users"`

	ID        uint64     `bun:",pk"`
	CreatedAt *time.Time `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
	Username  string     `bun:"username"`
	Password  string     `bun:"password"`
	Role      uint32     `bun:"role"`
}

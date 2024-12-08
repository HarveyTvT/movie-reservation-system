package model

import (
	"time"

	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
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

func (u *User) ToPB() *movie_reservation.User {
	if u == nil {
		return nil
	}
	return &movie_reservation.User{
		Username: u.Username,
		Role:     movie_reservation.User_Role(u.Role),
	}
}

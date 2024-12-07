package model

import "github.com/uptrace/bun"

type MovieGenre struct {
	bun.BaseModel `bun:"movie_genres,alias:movie_genres"`

	MovieID uint64 `bun:"movie_id"`
	Genre   string `bun:"genre"`
}

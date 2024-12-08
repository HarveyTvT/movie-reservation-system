package model

import "github.com/uptrace/bun"

type Genre struct {
	bun.BaseModel `bun:"genres,alias:genres"`

	Name string `bun:"name,pk"`
}

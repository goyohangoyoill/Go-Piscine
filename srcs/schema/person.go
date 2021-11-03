package schema

import (
	"time"
)

type Person struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Course    int       `json:"course"`
	Point     uint      `json:"point"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

package postgres

import (
	"time"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

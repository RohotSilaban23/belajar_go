package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Nama      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"created_at"`
}

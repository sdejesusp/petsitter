package models

import "time"

type User struct {
	FullName  string
	Email     string
	Password  string
	// Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

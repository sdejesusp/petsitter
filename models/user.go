package models

import "time"

type User struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	FullName  string   `json:"full_name"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles" gorm:"serializer:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

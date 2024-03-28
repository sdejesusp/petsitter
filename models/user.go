package models

import "time"

type User struct {
	ID        uint64   `json:"id" gorm:"primaryKey"`
	FullName  string   `json:"fullName"`
	Email     string   `json:"email" gorm:"uniqueIndex"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles" gorm:"serializer:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

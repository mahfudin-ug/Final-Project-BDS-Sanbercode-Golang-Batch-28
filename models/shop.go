package models

import "time"

type (
	Shop struct {
		ID        uint      `json:"id" gorm:"primary_key"`
		Name      string    `json:"name" gorm:"not null"`
		Bank      string    `json:"bank"`
		Phone     string    `json:"phone"`
		UserID    uint      `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		User      User      `json:"-"`
	}
)

package models

import "time"

type (
	Address struct {
		ID        uint      `json:"id" gorm:"primary_key"`
		Address   string    `json:"address" gorm:"not null"`
		City      string    `json:"city"`
		Country   string    `json:"country"`
		Province  string    `json:"province"`
		Zip       int       `json:"zip"`
		Latitude  string    `json:"latitude"`
		Longitude string    `json:"longitude"`
		UserID    uint      `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		User      User      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)

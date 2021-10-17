package models

import "time"

type (
	OrderProduct struct {
		ID        uint      `json:"id" gorm:"primary_key"`
		OrderID   uint      `json:"order_id"`
		ProductID uint      `json:"product_id"`
		Qty       int       `json:"qty"`
		Note      string    `json:"note"`
		Total     int       `json:"total"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Order     Order     `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		Product   Product   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)

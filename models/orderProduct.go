package models

import "time"

type (
	OrderProduct struct {
		ID        uint      `json:"id" gorm:"primary_key"`
		OrderID   uint      `json:"order_id"`
		ProductID uint      `json:"product_id"`
		Qty       uint      `json:"qty"`
		Note      string    `json:"note"`
		Total     uint      `json:"total"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Order     Order     `json:"-"`
		Product   Product   `json:"-"`
	}
)

package models

import "time"

type (
	Product struct {
		ID          uint      `json:"id" gorm:"primary_key"`
		Name        string    `json:"name" gorm:"not null"`
		Description string    `json:"desc"`
		Dimension   string    `json:"dimension"`
		PhotoPath   string    `json:"photo_path"`
		Price       uint      `json:"price"`
		Weight      uint      `json:"weight"`
		CategoryID  uint      `json:"category_id"`
		ShopID      uint      `json:"shop_id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Category    Category  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Shop        Shop      `json:"-"`
	}
)

package models

import "time"

type (
	Category struct {
		ID               uint      `json:"id" gorm:"primary_key"`
		Name             string    `json:"name" gorm:"not null"`
		Description      string    `json:"description"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		ParentCategoryID uint      `json:"parent_category_id"`
	}
)

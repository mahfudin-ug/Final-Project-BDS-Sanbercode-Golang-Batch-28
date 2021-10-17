package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type (
	Product struct {
		ID          uint      `json:"id" gorm:"primary_key"`
		Name        string    `json:"name" gorm:"not null"`
		Description string    `json:"desc"`
		Stock       int       `json:"stock"`
		PhotoPath   string    `json:"photo_path"`
		Price       int       `json:"price"`
		Weight      int       `json:"weight"`
		CategoryID  uint      `json:"category_id"`
		ShopID      uint      `json:"shop_id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Category    Category  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Shop        Shop      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)

func (p *Product) SubtractStock(number int, db *gorm.DB) (*Product, error) {
	if int(p.Stock) < number {
		return p, errors.New("number must less than product stock")
	}
	p.Stock -= number

	err := db.Model(&p).Updates(p).Error
	if err != nil {
		return &Product{}, err
	}

	return p, nil
}

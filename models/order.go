package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Order struct {
		ID           uint           `json:"id" gorm:"primary_key"`
		Status       string         `json:"status"`
		Payment      string         `json:"payment"`
		Courier      string         `json:"courier"`
		Total        int            `json:"total"`
		UserID       uint           `json:"user_id"`
		Resi         string         `json:"resi"`
		PaidAt       time.Time      `json:"paid_at" gorm:"type:TIMESTAMP NULL"`
		SentAt       time.Time      `json:"sent_at" gorm:"type:TIMESTAMP NULL"`
		CompleteAt   time.Time      `json:"complete_at" gorm:"type:TIMESTAMP NULL"`
		CreatedAt    time.Time      `json:"created_at"`
		UpdatedAt    time.Time      `json:"updated_at"`
		User         User           `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		OrderProduct []OrderProduct `json:"-"`
	}
)

const (
	OrderStatusInitial  = "INITIAL"
	OrderStatusPaid     = "PAID"
	OrderStatusSent     = "SENT"
	OrderStatusComplete = "COMPLETED"
)

func (o *Order) RecalculateOrder(db *gorm.DB) (*Order, error) {
	// Get all orderproduct
	var orderProducts []OrderProduct
	if err := db.Where("order_id=?", o.ID).Find(&orderProducts).Error; err != nil {
		return &Order{}, err
	}

	var grantTotal int
	for _, op := range orderProducts {
		grantTotal += op.Total
	}

	o.Total = grantTotal

	err := db.Model(&o).Updates(o).Error
	if err != nil {
		return &Order{}, err
	}

	return o, nil
}

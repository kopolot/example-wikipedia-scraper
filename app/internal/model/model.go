package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	ID        uint           `json:"id" gorm:"primaryKey"`
}

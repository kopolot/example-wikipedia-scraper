package model

import (
	modelTypes "example-wikipedia-scraper/internal/types/model"
)

type OrderItem struct {
	Order   Order   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Product Product `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Model
	OrderID   uint                     `json:"order_id" gorm:"not null;column:order_id;index:idx_order_item_order"`
	Quantity  uint                     `json:"quantity" gorm:"not null;column:quantity"`
	Price     float64                  `json:"price" gorm:"not null;column:price"`
	ItemType  modelTypes.OrderItemType `json:"item_type" gorm:"not null;column:item_type"`
	ProductID uint                     `json:"product_id" gorm:"not null;column:product_id"`
}

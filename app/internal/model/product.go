package model

import "example-wikipedia-scraper/internal/types/model"

type Product struct {
	Model
	Name           string            `json:"name" gorm:"not null;column:name"`
	Description    string            `json:"description" gorm:"column:description"`
	Price          float64           `json:"price" gorm:"not null;column:price"`
	Available      bool              `json:"available" gorm:"not null;column:available;default:true"`
	AvailableStock uint              `json:"availableStock" gorm:"not null;column:available_stock;default:0"`
	Type           model.ProductType `json:"type" gorm:"not null;column:type"`
}

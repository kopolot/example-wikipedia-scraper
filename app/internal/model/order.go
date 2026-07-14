package model

import (
	modelTypes "example-wikipedia-scraper/internal/types/model"
)

type Order struct {
	User User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Model
	Items         []*OrderItem           `json:"items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ClientIP      string                 `json:"client_ip" gorm:"not null;column:client_ip"`
	Status        modelTypes.OrderStatus `json:"status" gorm:"not null;column:status;index:idx_order_status"`
	InvoiceNumber *string                `json:"invoice_number" gorm:"column:invoice_number;uniqueIndex:idx_order_invoice_number"`
	Currency      string                 `json:"currency" gorm:"not null;column:currency;default:'EUR'"`
	PaymentID     string                 `json:"payment_id" gorm:"column:payment_id;index:idx_order_payment_id"`
	PaymentMethod string                 `json:"payment_method" gorm:"column:payment_method"`
	TotalAmount   float64                `json:"total_amount" gorm:"not null;column:total_amount"`
	UserID        uint                   `json:"user_id" gorm:"not null;column:user_id;index:idx_order_user"`
}

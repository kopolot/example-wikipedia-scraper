package dto

type OrderCreateDto struct {
	UserID        uint                  `json:"user_id" validate:"omitempty"`
	ClientIP      string                `json:"client_ip" validate:"omitempty,ip"`
	Currency      string                `json:"currency" validate:"omitempty,len=3"`
	PaymentMethod string                `json:"payment_method" validate:"required"`
	TotalAmount   float64               `json:"total_amount" validate:"required,gt=0"`
	WithInvoice   bool                  `json:"with_invoice"`
	Items         []*OrderItemCreateDto `json:"items" validate:"required,dive"`
}

type OrderItemCreateDto struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  uint `json:"quantity" validate:"required,gt=0"`
}

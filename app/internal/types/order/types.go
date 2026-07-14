package order

import "example-wikipedia-scraper/internal/types/model"

type OrderCreated struct {
	Status      model.OrderStatus `json:"status"`
	PaymentId   string            `json:"payment_id,omitempty"`
	RedirectURL string            `json:"redirect_url,omitempty"`
	ID          uint              `json:"id"`
}

type OrderPaymentNotification struct {
	ServiceIp     string            `json:"service_ip"`
	Body          string            `json:"body"`
	PaymentMethod string            `json:"payment_method"`
	Headers       map[string]string `json:"headers"`
}

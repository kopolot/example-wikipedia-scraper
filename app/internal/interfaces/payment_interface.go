package interfaces

import (
	orderTypes "example-wikipedia-scraper/internal/types/order"
	paymentTypes "example-wikipedia-scraper/internal/types/payment"
)

type PaymentInterface interface {
	CreatePayment(paymentRequest *paymentTypes.PaymentCreateRequest) (*paymentTypes.PaymentCreatedResponse, error)
	ValidatePaymentResponse(response *orderTypes.OrderPaymentNotification) (*paymentTypes.PaymentValidationResponse, error)
}

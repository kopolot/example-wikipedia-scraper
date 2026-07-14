package service

import (
	orderTypes "example-wikipedia-scraper/internal/types/order"
	paymentTypes "example-wikipedia-scraper/internal/types/payment"
	"fmt"
)

var (
	ErrInvalidPaymentMethod   = fmt.Errorf("invalid payment method")
	ErrDisabledPaymentMethod  = fmt.Errorf("payment method is disabled")
	ErrInvalidPaymentResponse = fmt.Errorf("invalid payment response")
)

type PaymentServiceInterface interface {
	CreatePayment(request *paymentTypes.PaymentCreateRequest, paymentMethod string) (*paymentTypes.PaymentCreatedResponse, error)
	GetAvailablePaymentMethods() []string
	HandlePaymentNotification(request *orderTypes.OrderPaymentNotification) (*paymentTypes.PaymentValidationResponse, error)
}

package registry

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/internal/service/payment/payments"
	"sync"
)

type PaymentFactory func() interfaces.PaymentInterface

var (
	paymentRegistry  map[string]PaymentFactory
	initOnceRegistry sync.Once
)

func NewPaymentRegistry(cfg config.ConfigInterface) map[string]PaymentFactory {
	initOnceRegistry.Do(func() {
		paymentRegistry = map[string]PaymentFactory{
			"example": func() interfaces.PaymentInterface {
				return payments.NewExamplePayment(cfg)
			},
		}
	})
	return paymentRegistry
}

func GetPaymentRegistry() map[string]PaymentFactory {
	return paymentRegistry
}

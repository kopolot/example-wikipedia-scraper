package payment

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/registry"
	orderTypes "example-wikipedia-scraper/internal/types/order"
	paymentTypes "example-wikipedia-scraper/internal/types/payment"
	"slices"
)

type PaymentService struct {
	registry map[string]registry.PaymentFactory
	cfg      config.ConfigInterface
}

func NewPaymentService(cfg config.ConfigInterface) *PaymentService {
	registry := registry.NewPaymentRegistry(cfg)
	return &PaymentService{
		registry: registry,
		cfg:      cfg,
	}
}

func (s *PaymentService) CreatePayment(paymentRequest *paymentTypes.PaymentCreateRequest, paymentMethod string) (*paymentTypes.PaymentCreatedResponse, error) {
	paymentProcessor, err := s.getPaymentProcessor(paymentMethod)
	if err != nil {
		return nil, err
	}
	response, err := paymentProcessor.CreatePayment(paymentRequest)
	return response, err
}

func (s *PaymentService) getPaymentProcessor(paymentMethod string) (interfaces.PaymentInterface, error) {
	paymentsConfig := s.cfg.GetPaymentMethodsConfig()
	paymentCfgIndex := slices.IndexFunc(paymentsConfig, func(v *config.PaymentMethodConfig) bool { return v.Name == paymentMethod })
	if paymentCfgIndex == -1 {
		return nil, service.ErrInvalidPaymentMethod
	}
	paymentCfg := paymentsConfig[paymentCfgIndex]
	if !paymentCfg.Enabled {
		return nil, service.ErrDisabledPaymentMethod
	}
	factory, exists := s.registry[paymentMethod]
	if !exists {
		return nil, service.ErrInvalidPaymentMethod
	}
	return factory(), nil
}

func (s *PaymentService) GetAvailablePaymentMethods() []string {
	paymentsConfig := s.cfg.GetPaymentMethodsConfig()
	availableMethods := make([]string, 0)
	for _, paymentCfg := range paymentsConfig {
		if paymentCfg.Enabled {
			availableMethods = append(availableMethods, paymentCfg.Name)
		}
	}
	return availableMethods
}

func (s *PaymentService) HandlePaymentNotification(request *orderTypes.OrderPaymentNotification) (*paymentTypes.PaymentValidationResponse, error) {
	paymentProcessor, err := s.getPaymentProcessor(request.PaymentMethod)
	if err != nil {
		return nil, err
	}
	paymentResponse, err := paymentProcessor.ValidatePaymentResponse(request)
	return paymentResponse, err
}

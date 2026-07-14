package payments

import (
	"example-wikipedia-scraper/internal/config"
	orderTypes "example-wikipedia-scraper/internal/types/order"
	paymentTypes "example-wikipedia-scraper/internal/types/payment"
	"encoding/json"
	"fmt"
	"time"
)

type ExamplePayment struct {
	cfg config.ConfigInterface
}

func NewExamplePayment(cfg config.ConfigInterface) *ExamplePayment {
	return &ExamplePayment{
		cfg: cfg,
	}
}

func (p *ExamplePayment) CreatePayment(paymentRequest *paymentTypes.PaymentCreateRequest) (*paymentTypes.PaymentCreatedResponse, error) {
	paymentID := fmt.Sprintf("pay_%d", time.Now().UnixNano())
	redirectURL := fmt.Sprintf(p.cfg.GetApiConfig().PublicHost+"order/example_payment/%s", paymentID)
	return &paymentTypes.PaymentCreatedResponse{
		Status:      paymentTypes.PaymentStatusPending,
		PaymentID:   paymentID,
		RedirectURL: redirectURL,
	}, nil
}

func (p *ExamplePayment) ValidatePaymentResponse(request *orderTypes.OrderPaymentNotification) (*paymentTypes.PaymentValidationResponse, error) {
	var notification struct {
		PaymentID string `json:"payment_id"`
		Action    string `json:"action"`
	}
	err := json.Unmarshal([]byte(request.Body), &notification)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment notification: %w", err)
	}
	result := &paymentTypes.PaymentValidationResponse{
		PaymentID:     notification.PaymentID,
		PaymentMethod: "example",
	}
	if notification.Action == "accept" {
		result.Status = paymentTypes.PaymentStatusCompleted
	}
	return result, nil
}

package payment

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

// for payment created response from operator
type PaymentCreatedResponse struct {
	Status      PaymentStatus `json:"status"`
	PaymentID   string        `json:"payment_id,omitempty"`
	RedirectURL string        `json:"redirect_url,omitempty"`
}

// for payment status returned from operator
type PaymentValidationResponse struct {
	Status        PaymentStatus `json:"status"`
	PaymentID     string        `json:"payment_id,omitempty"`
	OrderID       string        `json:"order_id,omitempty"`
	PaymentMethod string        `json:"payment_method,omitempty"`
}

// for payment creation request to operator
type PaymentCreateRequest struct {
	Items        []*PaymentItems `json:"items" validate:"required,dive"`
	OrderID      string          `json:"order_id" validate:"required"`
	ReturnURL    string          `json:"return_url" validate:"required,url"`
	NotifyURL    string          `json:"notify_url,omitempty" validate:"omitempty,url"`
	CurrencyCode string          `json:"currency_code" validate:"required,len=3"`
	Value        float64         `json:"value" validate:"required,gt=0"`
}

type PaymentItems struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description,omitempty"`
	Quantity    int     `json:"quantity" validate:"required,gt=0"`
	UnitPrice   float64 `json:"unit_price" validate:"required,gt=0"`
}

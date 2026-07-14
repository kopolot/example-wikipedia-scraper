package service

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
	modelTypes "example-wikipedia-scraper/internal/types/model"
	orderTypes "example-wikipedia-scraper/internal/types/order"
	paymentTypes "example-wikipedia-scraper/internal/types/payment"
	pkgRepository "example-wikipedia-scraper/pkg/repository"
)

type OrderWithItems struct {
	*model.Order
	OrderItems []*model.OrderItem
}

type OrderServiceInterface interface {
	CreateOrder(orderCreateDto *dto.OrderCreateDto) (*OrderWithItems, error)
	GetOrderItemType(product *model.Product) modelTypes.OrderItemType
	ProcessPayment(orderWithItems *OrderWithItems) (*paymentTypes.PaymentCreatedResponse, error)
	GetAvailablePaymentMethods() []string
	HandlePaymentNotification(request *orderTypes.OrderPaymentNotification) error
	GetOrderItemsByOrderID(orderID uint) []*model.OrderItem
	GetOrderByID(orderID uint) *model.Order
	GetProductRepo() pkgRepository.RepositoryInterface[*model.Product]
}

package repository

import (
	"example-wikipedia-scraper/internal/model"

	modelTypes "example-wikipedia-scraper/internal/types/model"
	pkgRepository "example-wikipedia-scraper/pkg/repository"
)

type OrderRepositoryInterface interface {
	pkgRepository.RepositoryInterface[*model.Order]
	UpdatePaymentID(orderID uint, paymentID string) error
	UpdateStatusByPaymentID(paymentID string, status modelTypes.OrderStatus) error
	UpdateStatus(orderID uint, status modelTypes.OrderStatus) error
}

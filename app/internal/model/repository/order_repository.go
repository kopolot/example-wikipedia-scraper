package repository

import (
	"example-wikipedia-scraper/internal/model"
	modelTypes "example-wikipedia-scraper/internal/types/model"
	pkgRepository "example-wikipedia-scraper/pkg/repository"
)

type OrderRepository struct {
	pkgRepository.GenericRepository[*model.Order]
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		GenericRepository: *pkgRepository.NewGenericRepository[*model.Order](),
	}
}

func (r *OrderRepository) UpdatePaymentID(orderID uint, paymentID string) error {
	qb := r.GetQueryBuilder()
	return qb.Model(&model.Order{}).Where("id = ?", orderID).Update("payment_id", paymentID)
}

func (r *OrderRepository) UpdateStatusByPaymentID(paymentID string, status modelTypes.OrderStatus) error {
	qb := r.GetQueryBuilder()
	return qb.Model(&model.Order{}).Where("payment_id = ?", paymentID).Update("status", status)
}

func (r *OrderRepository) UpdateStatus(orderID uint, status modelTypes.OrderStatus) error {
	qb := r.GetQueryBuilder()
	return qb.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status)
}

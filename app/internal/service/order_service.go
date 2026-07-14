package service

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
	queueInterfaces "example-wikipedia-scraper/internal/interfaces/queue"
	"example-wikipedia-scraper/internal/interfaces/repository"
	serviceInterfaces "example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/queue"
	modelTypes "example-wikipedia-scraper/internal/types/model"
	orderTypes "example-wikipedia-scraper/internal/types/order"
	paymentTypes "example-wikipedia-scraper/internal/types/payment"
	pkgDb "example-wikipedia-scraper/pkg/db"
	pkgRepository "example-wikipedia-scraper/pkg/repository"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/google/uuid"
)

var (
	ErrProductNotAvailable                = errors.New("product not available")
	ErrProductAvailableStockNotSufficient = errors.New("product available stock not sufficient")
)

type OrderService struct {
	orderRepo      repository.OrderRepositoryInterface
	orderItemRepo  pkgRepository.RepositoryInterface[*model.OrderItem]
	productRepo    pkgRepository.RepositoryInterface[*model.Product]
	paymentService serviceInterfaces.PaymentServiceInterface
	cfg            config.ConfigInterface
}

func NewOrderService(orderRepo repository.OrderRepositoryInterface, orderItemRepo pkgRepository.RepositoryInterface[*model.OrderItem], productRepo pkgRepository.RepositoryInterface[*model.Product], paymentService serviceInterfaces.PaymentServiceInterface, cfg config.ConfigInterface) *OrderService {
	// logger implementation
	return &OrderService{
		orderRepo:      orderRepo,
		orderItemRepo:  orderItemRepo,
		productRepo:    productRepo,
		paymentService: paymentService,
		cfg:            cfg,
	}
}

func (s *OrderService) CreateOrder(orderCreateDto *dto.OrderCreateDto) (*serviceInterfaces.OrderWithItems, error) {
	order := s.createOrderModelFromDto(orderCreateDto)
	if orderCreateDto.WithInvoice {
		invoiceNumber := s.generateInvoiceNumber()
		order.InvoiceNumber = &invoiceNumber
	}

	ordersItems := make([]*model.OrderItem, 0, len(orderCreateDto.Items))
	err := s.orderRepo.GetQueryBuilder().Transaction(func(tx pkgDb.QueryBuilder) error {
		if err := tx.Create(order); err != nil {
			return err
		}
		for _, itemDto := range orderCreateDto.Items {
			orderItem, err := s.createAndSaveOrderItem(itemDto, order, tx)
			if err != nil {
				return err
			}
			ordersItems = append(ordersItems, orderItem)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &serviceInterfaces.OrderWithItems{
		Order:      order,
		OrderItems: ordersItems,
	}, nil
}

func (s *OrderService) createAndSaveOrderItem(itemDto *dto.OrderItemCreateDto, order *model.Order, tx pkgDb.QueryBuilder) (*model.OrderItem, error) {
	product := &model.Product{}
	if err := tx.LockForUpdate().First(product, itemDto.ProductID); err != nil {
		return nil, err
	}

	stockTracked, err := s.handleProductStock(product, itemDto.Quantity)
	if err != nil {
		return nil, err
	}
	if stockTracked {
		if err := tx.Save(product); err != nil {
			return nil, err
		}
	}

	orderItem := &model.OrderItem{
		OrderID:   order.ID,
		ProductID: itemDto.ProductID,
		Quantity:  itemDto.Quantity,
		Price:     product.Price,
		ItemType:  s.GetOrderItemType(product),
	}
	if err := tx.Create(orderItem); err != nil {
		return nil, err
	}
	return orderItem, nil
}

func (s *OrderService) createOrderModelFromDto(orderCreateDto *dto.OrderCreateDto) *model.Order {
	return &model.Order{
		UserID:        orderCreateDto.UserID,
		ClientIP:      orderCreateDto.ClientIP,
		Currency:      orderCreateDto.Currency,
		PaymentMethod: orderCreateDto.PaymentMethod,
		TotalAmount:   orderCreateDto.TotalAmount,
		Status:        modelTypes.OrderStatusPending,
	}
}

func (s *OrderService) handleProductStock(product *model.Product, quantity uint) (bool, error) {
	if !product.Available {
		return false, ErrProductNotAvailable
	}
	if product.AvailableStock == 0 {
		return false, nil
	}
	if product.AvailableStock < quantity {
		return false, ErrProductAvailableStockNotSufficient
	}
	product.AvailableStock -= quantity
	if product.AvailableStock == 0 {
		product.Available = false
	}
	return true, nil
}

func (s *OrderService) getOrderItemTotal(orderItem *model.OrderItem) float64 {
	return float64(orderItem.Quantity) * orderItem.Price
}

func (s *OrderService) ProcessPayment(orderWithItems *serviceInterfaces.OrderWithItems) (*paymentTypes.PaymentCreatedResponse, error) {
	if orderWithItems.ID == 0 {
		return nil, errors.New("order ID is required for payment processing")
	}
	paymentRequest := s.createPaymentRequest(orderWithItems)
	response, err := s.paymentService.CreatePayment(paymentRequest, orderWithItems.PaymentMethod)
	if err != nil {
		return nil, err
	}
	if response.PaymentID == "" || response.Status != paymentTypes.PaymentStatusPending {
		return nil, serviceInterfaces.ErrInvalidPaymentResponse
	}
	err = s.orderRepo.UpdatePaymentID(orderWithItems.ID, response.PaymentID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *OrderService) createPaymentRequest(orderWithItems *serviceInterfaces.OrderWithItems) *paymentTypes.PaymentCreateRequest {
	strId := strconv.Itoa(int(orderWithItems.ID))
	paymentRequest := &paymentTypes.PaymentCreateRequest{
		OrderID:      strId,
		ReturnURL:    s.cfg.GetApiConfig().PublicFrontendHost + "panel/order/" + strId,
		CurrencyCode: orderWithItems.Currency,
		Value:        orderWithItems.TotalAmount,
		Items:        []*paymentTypes.PaymentItems{},
	}
	for _, item := range orderWithItems.OrderItems {
		name := item.Product.Name
		if name == "" {
			product, err := s.productRepo.FindOneBy("id = ?", item.ProductID)
			if err == nil && product != nil {
				name = product.Name
			}
		}
		paymentItem := &paymentTypes.PaymentItems{
			Name:      name,
			Quantity:  int(item.Quantity),
			UnitPrice: item.Price,
		}
		paymentRequest.Items = append(paymentRequest.Items, paymentItem)
	}
	return paymentRequest
}

func (s *OrderService) generateInvoiceNumber() string {
	// Implement your invoice number generation logic here
	return uuid.NewString()
}

func (s *OrderService) GetOrderItemType(product *model.Product) modelTypes.OrderItemType {
	var itemType modelTypes.OrderItemType
	switch product.Type {
	case modelTypes.ProductTypeSubscription:
		itemType = modelTypes.OrderItemTypeSubscription
	}
	return itemType
}

func (s *OrderService) GetAvailablePaymentMethods() []string {
	return s.paymentService.GetAvailablePaymentMethods()
}

func (s *OrderService) HandlePaymentNotification(request *orderTypes.OrderPaymentNotification) error {
	resp, err := s.paymentService.HandlePaymentNotification(request)
	if err == nil && resp.Status == paymentTypes.PaymentStatusCompleted {
		qb := s.orderRepo.GetQueryBuilder()
		if resp.OrderID == "" {
			order := &model.Order{}
			err = qb.Select("id").Where("payment_method = ? and payment_id = ?", resp.PaymentMethod, resp.PaymentID).First(order)
			if err == nil {
				resp.OrderID = strconv.Itoa(int(order.ID))
			}
		}
		updateQB := qb.Table("orders").Where("id = ?", resp.OrderID)
		err = updateQB.Update("status", modelTypes.OrderStatusPaid)
		if err == nil && updateQB.RowsAffected() > 0 {
			orderId, _ := strconv.Atoi(resp.OrderID)
			payload, _ := json.Marshal(uint(orderId))
			err = queue.GetMessageQueueService().Publish(&queueInterfaces.Task{
				Type:    "order_payment_notfied",
				Payload: queueInterfaces.JSONString(payload),
			})
		}
	}
	return err
}

func (s *OrderService) GetOrderItemsByOrderID(orderID uint) []*model.OrderItem {
	var orderItems []*model.OrderItem
	qb := s.orderItemRepo.GetQueryBuilder()
	err := qb.Where("order_id = ?", orderID).Preload("Product").Find(&orderItems)
	if err != nil {
		return nil
	}
	return orderItems
}

func (s *OrderService) GetOrderByID(orderId uint) *model.Order {
	order, err := s.orderRepo.GetByIDWithPreloads(orderId, "User", "Items")
	if err != nil || order == nil || order.ID == 0 {
		return nil
	}
	return order
}

func (s *OrderService) GetProductRepo() pkgRepository.RepositoryInterface[*model.Product] {
	return s.productRepo
}

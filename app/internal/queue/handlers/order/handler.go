package order

import (
	"example-wikipedia-scraper/internal/interfaces"
	"example-wikipedia-scraper/internal/interfaces/queue"
	serviceInterfaces "example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/model"
	modelTypes "example-wikipedia-scraper/internal/types/model"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type OrderQueueHandler struct {
	orderService serviceInterfaces.OrderServiceInterface
	subService   serviceInterfaces.SubscriptionServiceInterface
	logger       interfaces.LoggerInterface
}

func NewOrderQueueHandler(orderService serviceInterfaces.OrderServiceInterface, subService serviceInterfaces.SubscriptionServiceInterface, logger interfaces.LoggerInterface) *OrderQueueHandler {
	return &OrderQueueHandler{
		orderService: orderService,
		subService:   subService,
		logger:       logger,
	}
}

func (h *OrderQueueHandler) HandleOrderPaymentNotfied(task *queue.Task) {
	logger := h.logger
	var orderId uint
	err := json.Unmarshal([]byte(task.Payload), &orderId)
	if err != nil {
		logger.Error("Failed to unmarshal order ID from task payload", "payload", task.Payload, "error", err.Error())
		panic(fmt.Sprintf("Failed to unmarshal order ID from task payload: %v, error: %s", task.Payload, err.Error()))
	}
	order := h.orderService.GetOrderByID(orderId)
	user := &order.User
	timeToAdd := time.Duration(0)
	var currentLvl int64
	for _, item := range order.Items {
		product, err := h.orderService.GetProductRepo().FindOneBy("id = ?", item.ProductID)
		if err != nil {
			logger.Error("Failed to retrieve product for order item in order_payment_notfied handler", "orderID", order.ID, "productID", item.ProductID, "error", err.Error())
			panic(fmt.Sprintf("Failed to retrieve product for order item in order_payment_notfied handler, order ID: %d, product ID: %d, error: %s", order.ID, item.ProductID, err.Error()))
		}
		switch product.Type {
		case modelTypes.ProductTypeSubscription:
			var tmp time.Duration
			tmp, currentLvl = h.handlePaymentNotifiedSubscription(order, currentLvl, item)
			timeToAdd += tmp
		default:
		}
	}
	if timeToAdd != time.Duration(0) {
		err := h.subService.AddSubscriptionTime(user, currentLvl, timeToAdd)
		if err != nil {
			logger.Error("Failed to add subscription time", "userID", user.ID, "error", err.Error())
			panic(fmt.Sprintf("Failed to add subscription time for user ID %d: %s", user.ID, err.Error()))
		}
	}
}

func (h *OrderQueueHandler) handlePaymentNotifiedSubscription(order *model.Order, currentLvl int64, item *model.OrderItem) (time.Duration, int64) {
	var subLvlProd model.SubscriptionLevelProduct
	err := h.subService.GetSubscriptionLevelProductRepo().GetQueryBuilder().PreloadAssociations().Where("product_id = ?", item.ProductID).First(&subLvlProd)
	productIdString := strconv.FormatUint(uint64(item.ProductID), 10)
	if err != nil {
		h.logger.Error("Invalid subscription product in order item", "orderID", order.ID, "productID", item.ProductID)
		panic(fmt.Sprintf("Invalid subscription product in order item, order ID: %d, product ID: %s", order.ID, productIdString))
	}
	if currentLvl != 0 && currentLvl != int64(subLvlProd.SubscriptionLevel.Level) {
		h.logger.Error("Multiple subscription products with different levels in the same order", "orderID", order.ID)
		panic(fmt.Sprintf("Multiple subscription products with different levels in the same order, order ID: %d", order.ID))
	}
	currentLvl = int64(subLvlProd.SubscriptionLevel.Level)
	days, err := strconv.Atoi(subLvlProd.Product.Description)
	if err != nil {
		h.logger.Error("Invalid subscription product description, expected number of days", "description", subLvlProd.Product.Description, "productID", item.ProductID)
		panic(fmt.Sprintf("Invalid subscription product description, expected number of days, got: %s, product ID: %s", subLvlProd.Product.Description, productIdString))
	}
	return time.Duration(days) * time.Hour * 24 * time.Duration(item.Quantity), currentLvl
}

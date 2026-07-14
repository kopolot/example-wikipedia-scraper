package model

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusRefunded  OrderStatus = "refunded"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusFailed    OrderStatus = "failed"
)

type OrderItemType string

const (
	OrderItemTypeSubscription OrderItemType = "subscription"
)

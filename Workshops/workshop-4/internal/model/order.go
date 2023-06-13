package model

import (
	"time"
)

type OrderState string

const (
	OrderStateCreate   = OrderState("create")   // Заказ создается
	OrderStateProcess  = OrderState("process")  // Поиск водителя
	OrderStateProgress = OrderState("progress") // Заказ выполняется
	OrderStateComplete = OrderState("complete") // Заказ выполнен
)

type OrderID uint64

type PipelineOrder struct {
	Order    Order
	ClientID ClientID
	Err      error
}

type Order struct {
	ID       OrderID
	ClientID ClientID
	DriverID DriverID
	Tracking []OrderTracking
}

type OrderTracking struct {
	State OrderState
	Start time.Time
	End   time.Time
}

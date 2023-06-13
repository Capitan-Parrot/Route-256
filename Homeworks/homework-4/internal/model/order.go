package model

import (
	"time"
)

type OrderState string

// Состояния заказа
const (
	OrderStateCreate   = OrderState("Создан")    // Заказ создается
	OrderStateProcess  = OrderState("Обработан") // Заказ обрабатывается
	OrderStateComplete = OrderState("Завершён")  // Заказ выполнен
)

type OrderID uint64
type WorkerID uint64

type PipelineOrder struct {
	Order Order
	Err   error
}

type Order struct {
	ID            OrderID
	WorkerID      WorkerID
	WarehouseID   WarehouseID   `json:"ID склада"`
	PickUpPointID PickUpPointID `json:"ID пункта выдачи"`
	Tracking      []OrderTracking
}

type OrderTracking struct {
	State OrderState
	Start time.Time
	End   time.Time
}

package generator

import (
	"homework/internal/config"
	"homework/internal/model"
)

// OrderIDs - генератор ID заказов
func OrderIDs() <-chan model.OrderID {
	result := make(chan model.OrderID, config.OrdersTotal)
	go func() {
		defer close(result)
		for i := 0; i < config.OrdersTotal; i++ {
			result <- model.OrderID(i)
		}
	}()
	return result
}

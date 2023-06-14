package process

import (
	"context"
	"homework/internal/model"
	orderCore "homework/internal/pkg/service/order"
	"time"
)

// Process обрабатывает заказ, инициализирует ID склада
func Process(server *orderCore.Server, order model.Order) (model.Order, error) {
	start := time.Now().UTC()
	time.Sleep(2 * time.Second)

	order.WarehouseID = model.WarehouseID(order.ID % 2)
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateProcess,
		Start: start,
		End:   time.Now().UTC(),
	})

	if err := server.Update(order); err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func Pipeline(ctx context.Context, server *orderCore.Server, orders <-chan model.PipelineOrder) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)
	go func() {
		defer close(outCh)
		for order := range orders {
			if order.Err != nil {
				select {
				case <-ctx.Done():
					return
				case outCh <- order:
				}
			}

			orderR, err := Process(server, order.Order)
			select {
			case <-ctx.Done():
				return
			case outCh <- model.PipelineOrder{
				Order: orderR,
				Err:   err,
			}:
			}
		}
	}()

	return outCh
}

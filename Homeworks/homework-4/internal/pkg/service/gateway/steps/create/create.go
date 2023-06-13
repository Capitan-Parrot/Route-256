package create

import (
	"context"
	"homework/internal/model"
	ordercore "homework/internal/pkg/service/order"
	"time"
)

// Create создаёт заказ, инициализирует ID заказа и ID обрабатывающего воркера
func Create(server *ordercore.Server,
	orderID model.OrderID, workerID model.WorkerID) (model.Order, error) {

	start := time.Now().UTC()
	order := model.Order{
		ID:       orderID,
		WorkerID: workerID,
		Tracking: []model.OrderTracking{{
			State: model.OrderStateCreate,
			Start: start,
			End:   time.Now().UTC(),
		}},
	}
	time.Sleep(1 * time.Second)

	if err := server.Create(order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func Pipeline(ctx context.Context, server *ordercore.Server, orderIDCh <-chan model.OrderID,
	workerID model.WorkerID) <-chan model.PipelineOrder {

	outCh := make(chan model.PipelineOrder)
	go func() {
		defer close(outCh)
		for orderID := range orderIDCh {
			orderR, err := Create(server, orderID, workerID)
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

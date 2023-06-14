package complete

import (
	"context"
	"fmt"
	"homework/internal/model"
	ordercore "homework/internal/pkg/service/order"
	"log"
	"time"
)

// Complete создаёт заказ, инициализирует ID заказа и ID обрабатывающего воркера
func Complete(server *ordercore.Server, order model.Order) (model.Order, error) {
	start := time.Now().UTC()
	time.Sleep(1 * time.Second)

	order.PickUpPointID = model.PickUpPointID(order.ID) + model.PickUpPointID(order.WarehouseID)
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateComplete,
		Start: start,
		End:   time.Now().UTC(),
	})

	if err := server.Update(order); err != nil {
		log.Println(err)
		return model.Order{}, err
	}

	fmt.Printf("OrderId: [%d], WorkerID: [%d], Warehouse: [%d], PickUpPointId: [%d], Tracking: [%v]\n",
		order.ID, order.WorkerID, order.WarehouseID, order.PickUpPointID, order.Tracking)

	return order, nil
}

func Pipeline(ctx context.Context, server *ordercore.Server, orders <-chan model.PipelineOrder) <-chan model.PipelineOrder {
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

			orderR, err := Complete(server, order.Order)
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

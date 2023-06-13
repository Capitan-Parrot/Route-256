package progress

import (
	"context"
	"time"

	"workshop.4.2/internal/model"
	orderclient "workshop.4.2/internal/pkg/service/gateway/client/order"
)

func New(order *orderclient.Client) *Implementation {
	return &Implementation{
		order: order,
	}
}

type Implementation struct {
	order *orderclient.Client
}

func (i *Implementation) Progress(ctx context.Context, order model.Order) (model.Order, error) {
	start := time.Now().UTC()

	time.Sleep(3 * time.Second)

	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateProgress,
		Start: start,
		End:   time.Now().UTC(),
	})

	if err := i.order.Update(ctx, order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (i *Implementation) Pipeline(ctx context.Context, orders <-chan model.PipelineOrder) <-chan model.PipelineOrder {
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

			orderR, err := i.Progress(ctx, order.Order)
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

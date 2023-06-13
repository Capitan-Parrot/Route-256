package process

import (
	"context"
	"time"

	"workshop.4.2/internal/model"
	driverclient "workshop.4.2/internal/pkg/service/gateway/client/driver"
	orderclient "workshop.4.2/internal/pkg/service/gateway/client/order"
)

func New(order *orderclient.Client, driver *driverclient.Client) *Implementation {
	return &Implementation{
		driver: driver,
		order:  order,
	}
}

type Implementation struct {
	driver *driverclient.Client
	order  *orderclient.Client
}

func (i *Implementation) Process(ctx context.Context, order model.Order) (model.Order, error) {
	start := time.Now().UTC()

	time.Sleep(1 * time.Second)

	driver, err := i.driver.Free(ctx)
	if err != nil {
		return model.Order{}, err
	}

	order.DriverID = driver.ID
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateProcess,
		Start: start,
		End:   time.Now().UTC(),
	})

	if err = i.order.Update(ctx, order); err != nil {
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

			orderR, err := i.Process(ctx, order.Order)
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

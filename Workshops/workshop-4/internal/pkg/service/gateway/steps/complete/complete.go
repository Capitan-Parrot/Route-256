package complete

import (
	"context"
	"fmt"
	"log"
	"sync"

	"workshop.4.2/internal/model"
	driverclient "workshop.4.2/internal/pkg/service/gateway/client/driver"
	orderclient "workshop.4.2/internal/pkg/service/gateway/client/order"
	userclient "workshop.4.2/internal/pkg/service/gateway/client/user"
)

func New(user *userclient.Client, driver *driverclient.Client, order *orderclient.Client) *Implementation {
	return &Implementation{
		user:   user,
		driver: driver,
		order:  order,
	}
}

type Implementation struct {
	user   *userclient.Client
	driver *driverclient.Client
	order  *orderclient.Client
}

func (i *Implementation) Complete(ctx context.Context, order model.Order) (model.Order, error) {
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateComplete,
	})

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		if err := i.order.Update(ctx, order); err != nil {
			log.Println(err)
		}
	}()

	var client model.Client
	go func() {
		defer wg.Done()

		var err error
		client, err = i.user.Get(ctx, order.ClientID)
		if err != nil {
			log.Println(err)
		}
	}()

	var driver model.Driver
	go func() {
		defer wg.Done()

		var err error
		driver, err = i.driver.Get(ctx, order.DriverID)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	fmt.Printf("Driver: [%s], client: [%s], car: [%s], tracking: [%v]", driver.Name, client.Name, driver.Car, order.Tracking)

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

			orderR, err := i.Complete(ctx, order.Order)
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

package create

import (
	"context"
	"time"

	"workshop.4.2/internal/model"
	orderclient "workshop.4.2/internal/pkg/service/gateway/client/order"
	userclient "workshop.4.2/internal/pkg/service/gateway/client/user"
)

func New(user *userclient.Client, order *orderclient.Client, ids <-chan model.OrderID) *Implementation {
	return &Implementation{
		user:  user,
		order: order,
		ids:   ids,
	}
}

type Implementation struct {
	user  *userclient.Client
	order *orderclient.Client
	ids   <-chan model.OrderID
}

func (i *Implementation) Create(ctx context.Context, clientID model.ClientID) (model.Order, error) {
	start := time.Now().UTC()

	// Проверяем, что пользователь существует
	if _, err := i.user.Get(ctx, clientID); err != nil {
		return model.Order{}, err
	}

	orderID := <-i.ids
	order := model.Order{
		ID:       orderID,
		ClientID: clientID,
		Tracking: []model.OrderTracking{{
			State: model.OrderStateCreate,
			Start: start,
			End:   time.Now().UTC(),
		}},
	}

	if err := i.order.Create(ctx, order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (i *Implementation) Pipeline(ctx context.Context, clientIDCh <-chan model.ClientID) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)
	go func() {
		defer close(outCh)
		for clientID := range clientIDCh {
			orderR, err := i.Create(ctx, clientID)
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

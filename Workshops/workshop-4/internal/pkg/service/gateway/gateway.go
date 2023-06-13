package gateway

import (
	"context"
	"log"
	"sync"

	"workshop.4.2/internal/model"
	orderclient "workshop.4.2/internal/pkg/service/gateway/client/order"
	userclient "workshop.4.2/internal/pkg/service/gateway/client/user"
	completestep "workshop.4.2/internal/pkg/service/gateway/steps/complete"
	createstep "workshop.4.2/internal/pkg/service/gateway/steps/create"
	processstep "workshop.4.2/internal/pkg/service/gateway/steps/process"
	progressstep "workshop.4.2/internal/pkg/service/gateway/steps/progress"
)

func New(create *createstep.Implementation, process *processstep.Implementation, progress *progressstep.Implementation, complete *completestep.Implementation) *Implementation {
	return &Implementation{
		user:     nil,
		order:    nil,
		create:   create,
		process:  process,
		progress: progress,
		complete: complete,
	}
}

type Implementation struct {
	user     *userclient.Client
	order    *orderclient.Client
	create   *createstep.Implementation
	process  *processstep.Implementation
	progress *progressstep.Implementation
	complete *completestep.Implementation
}

func (i *Implementation) Process(ctx context.Context, clientID model.ClientID) error {
	orderCreated, err := i.create.Create(ctx, clientID)
	if err != nil {
		return err
	}

	orderProcessed, err := i.process.Process(ctx, orderCreated)
	if err != nil {
		return err
	}

	orderProgressed, err := i.progress.Progress(ctx, orderProcessed)
	if err != nil {
		return err
	}

	if _, err = i.complete.Complete(ctx, orderProgressed); err != nil {
		return err
	}

	return nil
}

func (i *Implementation) Pipeline(ctx context.Context, clientID <-chan model.ClientID) {
	createCh := i.create.Pipeline(ctx, clientID)
	processCh := i.process.Pipeline(ctx, createCh)
	progressCh := i.progress.Pipeline(ctx, processCh)
	completeCh := i.complete.Pipeline(ctx, progressCh)

	for order := range completeCh {
		if order.Err != nil {
			log.Printf("error while processing order for clientID: [%d], err: [%v]", order.ClientID, order.Err)
		}
	}
}

func (i *Implementation) PipelineFan(ctx context.Context, clientID <-chan model.ClientID) {
	createCh := i.create.Pipeline(ctx, clientID)
	processCh := i.process.Pipeline(ctx, createCh)

	const limit = 5
	fanOutProgress := make([]<-chan model.PipelineOrder, limit)
	for it := 0; it < limit; it++ {
		fanOutProgress[it] = i.progress.Pipeline(ctx, processCh)
	}

	pipeline := i.complete.Pipeline(ctx, fanIn(ctx, fanOutProgress))
	for order := range pipeline {
		if order.Err != nil {
			log.Printf("error while processing order for clientID: [%d], err: [%v]", order.ClientID, order.Err)
		}
	}
}

func fanIn(ctx context.Context, chans []<-chan model.PipelineOrder) <-chan model.PipelineOrder {
	muliteplexed := make(chan model.PipelineOrder)

	var wg sync.WaitGroup
	for _, ch := range chans {
		wg.Add(1)

		go func(ch <-chan model.PipelineOrder) {
			defer wg.Done()
			for v := range ch {
				select {
				case <-ctx.Done():
					return
				case muliteplexed <- v:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(muliteplexed)
	}()

	return muliteplexed
}

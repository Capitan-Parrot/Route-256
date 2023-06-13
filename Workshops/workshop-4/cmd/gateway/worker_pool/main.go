package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"workshop.4.2/internal/model"
	"workshop.4.2/internal/pkg/generator"
	"workshop.4.2/internal/pkg/producer"
	"workshop.4.2/internal/pkg/service/gateway"
	driverclient "workshop.4.2/internal/pkg/service/gateway/client/driver"
	orderclient "workshop.4.2/internal/pkg/service/gateway/client/order"
	userclient "workshop.4.2/internal/pkg/service/gateway/client/user"
	completestep "workshop.4.2/internal/pkg/service/gateway/steps/complete"
	createstep "workshop.4.2/internal/pkg/service/gateway/steps/create"
	processstep "workshop.4.2/internal/pkg/service/gateway/steps/process"
	progressstep "workshop.4.2/internal/pkg/service/gateway/steps/progress"
)

const workerTotal = 5

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := userclient.New()
	order := orderclient.New()
	driver := driverclient.New()
	ids := generator.OrderIDs(ctx)

	create := createstep.New(user, order, ids)
	process := processstep.New(order, driver)
	progress := progressstep.New(order)
	complete := completestep.New(user, driver, order)

	server := gateway.New(create, process, progress, complete)

	errCh := make(chan error)
	orders := producer.Orders()

	var wg sync.WaitGroup
	start := time.Now().UTC()
	for i := 0; i < workerTotal; i++ {
		wg.Add(1)
		go worker(ctx, &wg, server, orders, errCh)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		log.Println(err)
	}

	fmt.Printf("Total duration: %f", time.Since(start).Seconds())
}

func worker(ctx context.Context, wg *sync.WaitGroup, server *gateway.Implementation, orders <-chan model.ClientID, errCh chan<- error) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		case id, ok := <-orders:
			if !ok {
				return
			}

			if err := server.Process(ctx, id); err != nil {
				errCh <- err
			}
		}
	}
}

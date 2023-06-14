package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"homework/internal/config"
	"homework/internal/model"
	"homework/internal/pkg/service/gateway/steps/complete"
	"homework/internal/pkg/service/gateway/steps/create"
	"homework/internal/pkg/service/gateway/steps/process"
	ordercore "homework/internal/pkg/service/order"
	"log"
	"sync"
)

// Pipeline простой шаблон
func Pipeline(ctx context.Context, server *ordercore.Server, orderID <-chan model.OrderID, workerID model.WorkerID) {
	createCh := create.Pipeline(ctx, server, orderID, workerID) // создание заказа
	processCh := process.Pipeline(ctx, server, createCh)        // процессинг заказа
	completeCh := complete.Pipeline(ctx, server, processCh)     // завершение заказа

	for order := range completeCh {
		if order.Err != nil {
			log.Printf("error while processing order for orderID: [%d], err: [%v]", order.Order.ID, order.Err)
			return
		}
		jsonResult, err := json.Marshal(order.Order)
		if err != nil {
			log.Printf("error while marshalling order for orderID: [%d], err: [%v]", order.Order.ID, err)
			return
		}
		fmt.Println(string(jsonResult))
	}
}

// Pipeline шаблон с FanIn/FanOut
func PipelineFan(ctx context.Context, server *ordercore.Server, orderID <-chan model.OrderID, workerID model.WorkerID) {
	createCh := create.Pipeline(ctx, server, orderID, workerID) // создание заказа

	fanOutProcess := make([]<-chan model.PipelineOrder, config.FunOutChanLimit) // разделение каналов
	for it := 0; it < config.FunOutChanLimit; it++ {
		fanOutProcess[it] = process.Pipeline(ctx, server, createCh) // процессинг заказа
	}

	pipeline := complete.Pipeline(ctx, server, fanIn(ctx, fanOutProcess)) // завершение заказа
	for order := range pipeline {
		if order.Err != nil {
			log.Printf("error while processing order for orderID: [%d], err: [%v]", order.Order.ID, order.Err)
			return
		}
		jsonResult, err := json.Marshal(order.Order)
		if err != nil {
			log.Printf("error while marshalling order for orderID: [%d], err: [%v]", order.Order.ID, err)
			return
		}
		fmt.Println(string(jsonResult))
	}
}

// fanIn сливает потоки в один
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

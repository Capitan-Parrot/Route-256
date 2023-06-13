package main

import (
	"context"
	"fmt"
	"homework/internal/config"
	"homework/internal/model"
	"homework/internal/pkg/generator"
	"homework/internal/pkg/service/gateway"
	ordercore "homework/internal/pkg/service/order"
	"sync"
	"time"
)

// пайплайн без fanIn/fanOut
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// генерация заказов
	orderIDs := generator.OrderIDs()
	// создание сервера с заказами
	server := ordercore.New()
	start := time.Now().UTC()
	var wg sync.WaitGroup
	for i := 0; i < config.WorkerTotal; i++ {
		wg.Add(1)
		go worker(ctx, &wg, server, orderIDs, model.WorkerID(i))
	}

	wg.Wait()

	fmt.Printf("Total duration: %f", time.Since(start).Seconds())
}

func worker(ctx context.Context, wg *sync.WaitGroup, server *ordercore.Server, orders <-chan model.OrderID, workerID model.WorkerID) {
	defer wg.Done()
	gateway.Pipeline(ctx, server, orders, workerID)
}

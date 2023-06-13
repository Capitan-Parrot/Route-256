package main

import (
	"context"
	"fmt"
	"time"

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

	start := time.Now().UTC()

	orders := producer.Orders()
	server.Pipeline(ctx, orders)

	fmt.Printf("Total duration: %f", time.Since(start).Seconds())
}

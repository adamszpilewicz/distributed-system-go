package main

import (
	"context"
	"fmt"
	"github.com/adamszpilewicz/distributed-systems/app/grades"
	"github.com/adamszpilewicz/distributed-systems/app/registry"
	"github.com/adamszpilewicz/distributed-systems/app/service"
	stlog "log"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("%v:%v", host, port)

	var r registry.Registration
	r.ServiceName = registry.GradingService
	r.ServiceURL = serviceAddress

	ctx, err := service.Start(context.Background(), host, port, r, grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("shutting down grading service")
}

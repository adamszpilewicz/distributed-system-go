package main

import (
	"context"
	"fmt"
	"github.com/adamszpilewicz/distributed-systems/app/grades"
	"github.com/adamszpilewicz/distributed-systems/app/log"
	"github.com/adamszpilewicz/distributed-systems/app/registry"
	"github.com/adamszpilewicz/distributed-systems/app/service"
	stlog "log"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.Registration
	r.ServiceName = registry.GradingService
	r.ServiceURL = serviceAddress
	r.RequiredServices = []registry.ServiceName{registry.LogService}
	r.ServiceUpdateURL = r.ServiceURL + "/services"

	ctx, err := service.Start(context.Background(), host, port, r, grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("\nlogging service found %v", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("shutting down grading service")
}

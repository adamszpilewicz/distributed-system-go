package main

import (
	"context"
	"fmt"
	"github.com/adamszpilewicz/distributed-systems/app/log"
	"github.com/adamszpilewicz/distributed-systems/app/registry"
	"github.com/adamszpilewicz/distributed-systems/app/service"

	stlog "log"
)

func main() {
	log.Run("./app.log")

	host, port := "localhost", "4000"
	var r registry.Registration
	r.ServiceName = registry.LogService
	r.ServiceURL = fmt.Sprintf("http://%v:%v", host, port)

	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("shutting down the log service")

}

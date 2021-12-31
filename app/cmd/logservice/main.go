package main

import (
	"context"
	"fmt"
	"github.com/adamszpilewicz/distributed-systems/app/log"
	"github.com/adamszpilewicz/distributed-systems/app/service"

	stlog "log"

)

func main() {
	log.Run("./app.log")

	host, port := "localhost", "4000"

	ctx, err := service.Start(
		context.Background(),
		"Log service",
		host,
		port,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatal(err)
	}

	<- ctx.Done()
	fmt.Println("shutting down the log service")

}





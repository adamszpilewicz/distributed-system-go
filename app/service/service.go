package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, serviceName, host, port string,
	registerHandlersFunc func()) (context.Context, error) {

	registerHandlersFunc()
	ctx = startService(ctx, serviceName, host, port)

	return ctx, nil
}

func startService(ctx context.Context, name string, host string,
	port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("server started: %v" +
			"\nif you want to shutdown the server then press any key", name)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}

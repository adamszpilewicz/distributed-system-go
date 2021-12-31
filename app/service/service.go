package service

import (
	"context"
	"fmt"
	"github.com/adamszpilewicz/distributed-systems/app/registry"
	"log"
	"net/http"
)

func Start(ctx context.Context, host, port string, reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {

	registerHandlersFunc()
	ctx = startService(ctx, string(reg.ServiceName), host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
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
		fmt.Printf("server started: %v at host %v"+
			"\nif you want to shutdown the server then press any key", name, host)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%v:%v", host, port))
		if err != nil {
			fmt.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}

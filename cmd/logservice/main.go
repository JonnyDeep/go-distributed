package main

import (
	"context"
	"fmt"
	stdlog "log"

	"huangliqun.github.com/log"
	"huangliqun.github.com/registry"
	"huangliqun.github.com/service"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: "Log Service",
		ServiceUrl:  serviceAddress,
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.Registerhandlers,
	)

	if err != nil {
		stdlog.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down log service.")

}

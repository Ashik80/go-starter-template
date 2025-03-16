package main

import (
	"context"
	"log"
	"net/http"

	"go-starter-template/pkg/bootstrap"
)

func main() {
	ctx := context.Background()
	app := bootstrap.Init(ctx)

	go func() {
		if err := app.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start application: %v", err)
		}
	}()

	app.GracefulShutdown(ctx)
}

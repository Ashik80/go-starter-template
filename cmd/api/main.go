package main

import (
	"context"
	"net/http"

	"go-starter-template/internal/bootstrap"
)

func main() {
	ctx := context.Background()
	app := bootstrap.Init(ctx)
	log := app.GetLogger()

	go func() {
		if err := app.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start application: %v", err)
		}
	}()

	app.GracefulShutdown(ctx)
}

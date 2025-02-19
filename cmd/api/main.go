package main

import (
	"context"
	"log"
	"net/http"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/handlers"
)

func main() {
	ctx := context.Background()

	app := app.Init(ctx)

	handlers.RegisterRoutes(app)

	go func() {
		if err := app.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start application: %v", err)
		}
	}()

	app.GracefulShutdown(ctx)
}

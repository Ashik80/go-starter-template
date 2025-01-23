package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-starter-template/ent"
	"go-starter-template/pkg/service"
	"go-starter-template/pkg/store"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type App struct {
	Router service.Router
	Orm    *ent.Client
	Store  *store.Store
	server *http.Server
}

func Init(ctx context.Context) *App {
	a := new(App)

	a.initRouterMux()
	a.initServer()
	a.initOrm()
	a.autoMigrateSchema(ctx)
	a.initStores()

	return a
}

func (a *App) initServer() {
	host := ""
	port := "8000"
	addr := fmt.Sprintf("%s:%s", host, port)
	server := http.Server{
		Addr:    addr,
		Handler: a.Router,
	}
	a.server = &server
}

func (a *App) initRouterMux() {
	a.Router = service.NewChiServerMux()
}

func (a *App) initOrm() {
	drv, err := entsql.Open(dialect.Postgres, "postgresql://postgres:postgres@localhost:5432/test_temp")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	client := ent.NewClient(ent.Driver(drv))
	a.Orm = client
}

func (a *App) autoMigrateSchema(ctx context.Context) {
	if err := a.Orm.Schema.Create(ctx); err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	}
}

func (a *App) initStores() {
	a.Store = store.NewDataStore(a.Orm)
}

func (a *App) Serve() error {
	log.Printf("server running on port %s\n", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) GracefulShutdown(ctx context.Context) {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	<-quitCh

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	log.Println("server shut down gracefully.")
}

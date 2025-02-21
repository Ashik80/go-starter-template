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
	"go-starter-template/pkg/config"
	"go-starter-template/pkg/service"
	"go-starter-template/pkg/store"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type App struct {
	Config           *config.Config
	Router           service.Router
	Orm              *ent.Client
	Store            *store.Store
	TemplateRenderer *service.TemplateRenderer
	server           *http.Server
	PasswordHasher   service.PasswordHasher
}

func Init(ctx context.Context) *App {
	a := new(App)

	a.initConfig()
	a.initOrm()
	// a.autoMigrateSchema(ctx)
	a.initRouterMux()
	a.initFileServer()
	a.initTemplatingEngine()
	a.initStores()
	a.initPasswordHasher()
	a.initApplicationServer()

	return a
}

func (a *App) initPasswordHasher() {
	a.PasswordHasher = service.NewBcryptPasswordHasher()
}

func (a *App) initConfig() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	a.Config = conf
}

func (a *App) initApplicationServer() {
	addr := fmt.Sprintf(":%s", a.Config.Port)
	server := http.Server{
		Addr:    addr,
		Handler: a.Router,
	}
	a.server = &server
}

func (a *App) initFileServer() {
	fs := http.FileServer(http.Dir("./web"))
	a.Router.Handle("/web/", http.StripPrefix("/web/", fs))
	log.Println("INFO: file server initialized in directory web/ directory")
}

func (a *App) initRouterMux() {
	a.Router = service.NewNetServerMux(a.Config)
	log.Println("INFO: router initialized")
}

func (a *App) initOrm() {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		a.Config.DatabaseConfig.User,
		a.Config.DatabaseConfig.Password,
		a.Config.DatabaseConfig.Host,
		a.Config.DatabaseConfig.Port,
		a.Config.DatabaseConfig.Name,
	)
	drv, err := entsql.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatalf("ERROR: failed to open database: %v", err)
	}
	client := ent.NewClient(ent.Driver(drv))
	a.Orm = client
	log.Println("INFO: ent orm initialized")
}

func (a *App) autoMigrateSchema(ctx context.Context) {
	if err := a.Orm.Schema.Create(ctx); err != nil {
		log.Fatalf("ERROR: failed to migrate schema: %v", err)
	}
	log.Println("INFO: auto migration initialized")
}

func (a *App) initStores() {
	a.Store = store.NewDataStore(a.Orm)
	log.Println("INFO: stores initialized")
}

func (a *App) initTemplatingEngine() {
	tr, err := service.NewTemplateRenderer(
		"web/templates/base.html",
		"web/templates/layouts",
		"web/templates/pages",
		"web/templates/partials",
	)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	a.TemplateRenderer = tr
	log.Println("INFO: templating engine initialized")
}

func (a *App) Serve() error {
	log.Printf("INFO: server running on port %s\n", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) GracefulShutdown(ctx context.Context) {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	<-quitCh

	log.Println("INFO: shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("ERROR: failed to shutdown server: %v", err)
	}

	log.Println("INFO: server shut down gracefully")
}

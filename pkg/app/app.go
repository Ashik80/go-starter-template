package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-starter-template/pkg/infrastructure"
	"go-starter-template/pkg/repository"
	"go-starter-template/pkg/service"

	_ "github.com/lib/pq"
)

type App struct {
	Config           *infrastructure.Config
	Router           infrastructure.Router
	DB               *sql.DB
	Repository       *repository.Repository
	Services         *service.Services
	TemplateRenderer infrastructure.TemplateRenderer
	server           *http.Server
	PasswordHasher   infrastructure.PasswordHasher
}

func Init(ctx context.Context) *App {
	a := new(App)

	a.initConfig()
	a.initDB()
	a.initRepositories()
	a.initServices()
	a.initRouterMux()
	a.initFileServer()
	a.initTemplatingEngine()
	a.initPasswordHasher()
	a.initApplicationServer()

	return a
}

func (a *App) initPasswordHasher() {
	a.PasswordHasher = infrastructure.NewBcryptPasswordHasher()
}

func (a *App) initConfig() {
	conf, err := infrastructure.NewConfig()
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
	a.Router = infrastructure.NewNetServerMux(a.Config)
	log.Println("INFO: router initialized")
}

func (a *App) initDB() {
	db, err := infrastructure.NewDatabaseConfig(a.Config)
	if err != nil {
		log.Fatalf("ERROR: failed to open database: %v", err)
	}
	a.DB = db
	log.Println("INFO: database initialized")
}

func (a *App) initRepositories() {
	a.Repository = repository.NewRepository(a.DB)
	log.Println("INFO: stores initialized")
}

func (a *App) initServices() {
	a.Services = service.NewServices(a.Repository)
	log.Println("INFO: services initialized")
}

func (a *App) initTemplatingEngine() {
	tr, err := infrastructure.NewTemplateRenderer(
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("ERROR: failed to shutdown server: %v", err)
	}

	log.Println("INFO: server shut down gracefully")
}

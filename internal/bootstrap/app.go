package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-starter-template/internal/infrastructure/config"
	"go-starter-template/internal/infrastructure/db/postgres"
	"go-starter-template/internal/infrastructure/factories"
	"go-starter-template/internal/infrastructure/middlewares"
	"go-starter-template/internal/interfaces/controllers"
	"go-starter-template/pkg/logger"
	"go-starter-template/pkg/renderer"
	"go-starter-template/pkg/router"

	_ "github.com/lib/pq"
)

type App struct {
	Config *config.Config
	Router router.Router
	DB     *sql.DB
	server *http.Server
	log    *logger.Logger
}

func Init(ctx context.Context) *App {
	a := new(App)

	// maintain the order to avoid nil pointer exception
	a.initLogger()
	a.initConfig()
	a.initDB()

	// if not using templ package enable this to use std templates
	// a.initTemplatingEngine()

	a.initRouterMux()
	a.initFileServer()
	a.initControllers()
	a.initApplicationServer()

	return a
}

func (a *App) GetLogger() *logger.Logger {
	return a.log
}

func (a *App) initControllers() {
	controllers.NewHealthController(a.Router, a.Config, a.DB)
	controllers.NewHomeController(a.Router)
	controllers.NewTodoController(
		a.Router,
		factories.NewTodoServiceWithPQRepository(a.DB),
		factories.NewSessionServiceWithPQRepository(a.DB),
		a.Config,
	)
	controllers.NewAuthController(
		a.Router,
		factories.NewUserServiceWithPQRepository(a.DB, a.log),
		a.Config,
	)
}

func (a *App) initLogger() {
	a.log = logger.NewLogger()
	a.log.Info("logger initialized")
}

func (a *App) initConfig() {
	conf, err := config.NewConfig()
	if err != nil {
		a.log.Fatal("%v", err)
	}
	a.Config = conf
	a.log.Info("configuration initialized")
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
	a.log.Info("file server initialized in directory web/ directory")
}

func (a *App) initRouterMux() {
	r := router.NewNetServerMux()

	// add middlewares
	r.Use(middlewares.Logger)
	csrfMiddleware := middlewares.CSRFMiddleware(a.Config.CSRFAuthKey)
	r.Use(csrfMiddleware)
	r.Use(middlewares.EnableCors(a.Config.AllowedOrigins))

	a.Router = r

	a.log.Info("router initialized")
}

func (a *App) initDB() {
	db, err := postgres.NewDatabaseConfig(a.Config)
	if err != nil {
		a.log.Fatal("failed to open database: %v", err)
	}
	a.DB = db
	a.log.Info("database initialized")
}

func (a *App) initTemplatingEngine() {
	err := renderer.InitBaseTemplate(a.log)
	if err != nil {
		a.log.Fatal("%v", err)
	}
	renderer.RegisterPageTemplates()
	a.log.Info("templates parsed")
}

func (a *App) Serve() error {
	a.log.Info("server running on port %s\n", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) GracefulShutdown(ctx context.Context) {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	<-quitCh

	a.log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.log.Fatal("failed to shutdown server: %v", err)
	}

	a.log.Info("server shut down gracefully")
}

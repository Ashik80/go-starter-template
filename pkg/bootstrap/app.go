package bootstrap

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

	"go-starter-template/pkg/application/services"
	"go-starter-template/pkg/infrastructure/config"
	"go-starter-template/pkg/infrastructure/db/postgres"
	"go-starter-template/pkg/infrastructure/renderer"
	"go-starter-template/pkg/infrastructure/router"
	"go-starter-template/pkg/infrastructure/security"
	"go-starter-template/pkg/interfaces/controllers"

	_ "github.com/lib/pq"
)

type App struct {
	Config *config.Config
	Router router.Router
	DB     *sql.DB
	// Repository *repository.Repository
	// Services   *service.Services
	server *http.Server
}

func Init(ctx context.Context) *App {
	a := new(App)

	// maintain the order to avoid nil pointer exception
	a.initConfig()
	a.initDB()
	a.initTemplatingEngine()
	a.initRouterMux()
	a.initFileServer()
	a.initControllers()
	a.initApplicationServer()

	return a
}

func (a *App) initControllers() {
	sessionRepo := postgres.NewPQSessionStore(a.DB)
	sessionService := services.NewSessionService(sessionRepo)

	controllers.NewHealthController(a.Router, a.Config, a.DB)
	controllers.NewHomeController(a.Router)

	todoRepo := postgres.NewPQTodoRepository(a.DB)
	todoService := services.NewTodoService(todoRepo)
	controllers.NewTodoController(a.Router, todoService, sessionService, a.Config)

	passwordHasher := security.NewBcryptPasswordHasher()
	userRepo := postgres.NewPQUserStore(a.DB)
	userService := services.NewUserService(userRepo, passwordHasher, sessionService)
	controllers.NewAuthController(a.Router, userService, a.Config)
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
	a.Router = router.NewNetServerMux(a.Config)
	log.Println("INFO: router initialized")
}

func (a *App) initDB() {
	db, err := postgres.NewDatabaseConfig(a.Config)
	if err != nil {
		log.Fatalf("ERROR: failed to open database: %v", err)
	}
	a.DB = db
	log.Println("INFO: database initialized")
}

func (a *App) initTemplatingEngine() {
	err := renderer.InitBaseTemplate()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	renderer.RegisterPageTemplates()
	log.Println("INFO: templates parsed")
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

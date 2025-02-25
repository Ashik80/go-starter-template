package handlers

import (
	"net/http"
	"runtime"
	"time"

	"go-starter-template/pkg/app"
	"go-starter-template/pkg/infrastructure"
)

type HealthHandler struct {
	Router infrastructure.Router
	app    *app.App
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	System    SystemInfo        `json:"system"`
}

type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}

func init() {
	Register(new(HealthHandler))
}

func (h *HealthHandler) Init(a *app.App) error {
	h.Router = a.Router
	h.app = a
	return nil
}

func (h *HealthHandler) Routes() {
	h.Router.Get("/health", h.Check)
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0", // TODO: Get from config
		Services: map[string]string{
			"database": "up",
		},
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			GOOS:         runtime.GOOS,
			GOARCH:       runtime.GOARCH,
			NumCPU:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
		},
	}

	if err := h.app.DB.PingContext(r.Context()); err != nil {
		status.Status = "unhealthy"
		status.Services["database"] = "down"
	}

	jsonResponse(w, http.StatusOK, status)
}

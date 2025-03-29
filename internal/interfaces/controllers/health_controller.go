package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"go-starter-template/internal/infrastructure/config"
	"go-starter-template/pkg/router"
)

type HealthHandler struct {
	conf *config.Config
	*sql.DB
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

func NewHealthController(r router.Router, conf *config.Config, sqlDB *sql.DB) {
	controller := &HealthHandler{
		conf: conf,
		DB:   sqlDB,
	}

	r.Get("/health", controller.Check)
}

func (hc *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   hc.conf.Version,
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

	if err := hc.DB.PingContext(r.Context()); err != nil {
		status.Status = "unhealthy"
		status.Services["database"] = "down"
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(status)
}

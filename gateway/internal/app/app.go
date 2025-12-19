package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
	} `yaml:"app"`
	HTTP struct {
		Port string `yaml:"port"`
	} `yaml:"http"`
	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`
	Upstreams struct {
		Order     string `yaml:"order"`
		Payment   string `yaml:"payment"`
		Analytics string `yaml:"analytics"`
		Menu      string `yaml:"menu"`
	} `yaml:"upstreams"`
}

type App struct {
	cfg    Config
	server *http.Server
}

func New(configPath string) *App {
	var cfg Config
	if configPath == "" {
		configPath = "/config/config.yaml"
	}

	if _, err := os.Stat(configPath); err == nil {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("failed to read config: %v", err)
		}
	} else {
		// defaults
		cfg.HTTP.Port = "8080"
		cfg.Logger.Level = "info"
		cfg.Upstreams.Order = "http://order-service:8080"
		cfg.Upstreams.Payment = "http://payment-service:8081"
		cfg.Upstreams.Analytics = "http://analytics-service:8083"
		cfg.Upstreams.Menu = "http://menu-service:8084"
		log.Warnf("config file not found, using defaults: %v", err)
	}

	initLogger(cfg.Logger.Level)

	return &App{cfg: cfg}
}

func (a *App) Start() {
	router := NewRouter(a.cfg)

	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", a.cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.WithField("port", a.cfg.HTTP.Port).Info("gateway starting")
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("gateway ListenAndServe: %v", err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.server == nil {
		return nil
	}
	return a.server.Shutdown(ctx)
}

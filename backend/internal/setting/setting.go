package setting

import (
	"crocus-calculator/internal/lib/env"
	"crocus-calculator/internal/lib/logger"
	cfg "crocus-calculator/internal/setting/config"
	"fmt"

	"go.uber.org/zap"
)

type App struct {
	cfg    cfg.Config
	logger *zap.SugaredLogger
}

func (a App) LoadConfig() {
	a.cfg = cfg.Config{
		AppEnv: env.GetString("APP_ENV", "development"),
		Server: cfg.ServerConfig{
			Addr: env.GetString("ADDR", ":8080"),
		},
		DB: cfg.DBConfig{
			Host:     env.GetString("DB_HOST", "localhost"),
			Port:     env.GetString("DB_PORT", "5432"),
			User:     env.GetString("DB_USER", "postgres"),
			Password: env.GetString("DB_PASSWORD", "postgres"),
			Name:     env.GetString("DB_NAME", "crocus"),
			SSLMode:  env.GetString("DB_SSLMODE", "disable"),
		},
	}
}

func (a *App) SetupLogger() {
	var err error
	a.logger, err = logger.Init(a.cfg.AppEnv)

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}

func (a App) InitApp() {
	a.LoadConfig()
	a.SetupLogger()
}

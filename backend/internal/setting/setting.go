package setting

import (
	"crocus-calculator/internal/api/auth"
	"crocus-calculator/internal/lib/env"
	"crocus-calculator/internal/lib/hasher"
	"crocus-calculator/internal/lib/logger"
	tokenservice "crocus-calculator/internal/lib/tokenService"
	"crocus-calculator/internal/lib/validator"
	cfg "crocus-calculator/internal/setting/config"
	"crocus-calculator/internal/storage"
	"crocus-calculator/internal/storage/repository"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type App struct {
	cfg    cfg.Config
	logger *zap.SugaredLogger
}

func (a *App) LoadConfig() {
	a.cfg = cfg.Config{
		AppEnv: env.GetString("APP_ENV", "development"),
		Server: cfg.ServerConfig{
			Addr: env.GetString("ADDR", ":8080"),
		},
		HTTP: cfg.HTTPConfig{
			AccessSecret:  env.GetString("HTTP_ACCESS_SECRET", "randomaccesssecret"),
			RefreshSecret: env.GetString("HTTP_REFRESH_SECRET", "randomrefreshsecret"),
			AccessTTL:     env.GetDuration("HTTP_ACCESS_TTL", 1*time.Hour),
			RefreshTTL:    env.GetDuration("HTTP_REFRESH_TTL", 30*24*time.Hour),
		},
		DB: cfg.DBConfig{
			Driver:   env.GetString("DB_DRIVER", "postgres"),
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

func (a App) Migrate() error {
	a.LoadConfig()

	dsn := a.cfg.DB.DSN()

	err := storage.Migrate(dsn)
	return err
}

func (a App) BootstrapDB() (*sqlx.DB, error) {
	a.LoadConfig()

	dsn := a.cfg.DB.DSN()
	db, err := storage.BootstrapDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to bootstrap database: %w", err)
	}

	return db, nil
}

func (a App) MountRouter(db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()

	hasher := hasher.BcryptHasher{}
	httpConfig := a.cfg.HTTP
	tknService := tokenservice.NewTokenService(httpConfig.AccessSecret, httpConfig.RefreshSecret, httpConfig.AccessTTL, httpConfig.RefreshTTL)
	authRepo := repository.NewAuthPostgresRepo(db)
	authController := auth.NewAuthController(auth.NewAuthService(authRepo, hasher, a.logger, tknService), a.logger, tknService)

	authController.RegisterRoutes(router)

	validator.InitValidator()
	return router
}

func (a App) StartServer(router *chi.Mux) {
	a.logger.Infof("Starting server on %s", a.cfg.Server.Addr)
	if err := http.ListenAndServe(a.cfg.Server.Addr, router); err != nil {
		a.logger.Fatalf("Failed to start server: %v", err)
	}
}

func (a App) InitApp() {
	a.LoadConfig()
	a.SetupLogger()
	db, err := a.BootstrapDB()
	if err != nil {
		a.logger.Fatalf("failed to bootstrap database: %v", err)
	}
	router := a.MountRouter(db)

	a.StartServer(router)
}

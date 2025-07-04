package setting

import (
	"crocus-calculator/internal/lib"
	cfg "crocus-calculator/internal/setting/config"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	cfg cfg.Config
}

func (a App) LoadConfig() {
	a.cfg = cfg.Config{
		AppEnv: lib.GetStringFromEnv("APP_ENV", "development"),
		Server: cfg.ServerConfig{
			Addr: lib.GetStringFromEnv("ADDR", ":8080"),
		},
	}
}

func (a App) MountRouter() *chi.Mux {
	r := chi.NewRouter()

	return r
}

func (a App) RunServer(r *chi.Mux) {
	http.ListenAndServe(a.cfg.Server.Addr, r)
}

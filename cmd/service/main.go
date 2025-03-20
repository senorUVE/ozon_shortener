package main

import (
	"log/slog"
	"net/http"
	"os"
	handlers "ozon_shortener/internal/api/url"
	"ozon_shortener/internal/config"
	"ozon_shortener/internal/middleware/validator"
	"ozon_shortener/internal/repository"
	"ozon_shortener/internal/repository/memory"
	"ozon_shortener/internal/services/url"

	_ "ozon_shortener/docs"

	_ "github.com/lib/pq"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
)

const (
	successExitCode = 0
	failExitCode    = 1
	letters         = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
)

// @title      URL shortener
// @version    1.0
// @description  URL shortener on go
// @host      localhost:8080
// @BasePath    /api
func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	v := validator.NewValidator(letters)
	cfg := config.Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("init config", "err", err)
		return failExitCode
	}

	db, err := sqlx.Open("postgres", cfg.DB.Dsn)
	if err != nil {
		slog.Error("init db", "err", err)
		return failExitCode
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("ping db", "err", err)
		return failExitCode
	}

	var dao repository.DAO
	switch cfg.App.StorageType {
	case "memory":
		memRepo := memory.NewMemoryRepository()
		dao = memory.NewMemoryDAO(memRepo)
	default:
		dao = repository.NewRepository(db)
	}
	urlSrv := url.New(dao, cfg.App.Url)

	handler := handlers.New(urlSrv, v)

	r := mux.NewRouter()

	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	r.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	apiRouter := r.PathPrefix("/api").Subrouter()

	urlRouter := apiRouter.PathPrefix("/url").Subrouter()

	urlRouter.HandleFunc("/generate", handler.CreateURL).Methods("POST")
	urlRouter.HandleFunc("/original", handler.GetOriginal).Methods("GET")

	urlRouter.HandleFunc("/{token}", handler.RedirectToOriginal).Methods("GET")

	slog.Info("start serve", slog.String("url", cfg.App.ListenAddr))
	if err := http.ListenAndServe(cfg.App.ListenAddr, r); err != nil {
		slog.Error("serve", "err", err)
		return failExitCode
	}

	return successExitCode
}

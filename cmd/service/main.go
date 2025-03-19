package main

import (
	"log/slog"
	"net/http"
	"os"
	adapters "ozon_shortener/internal/adapters/url"
	handlers "ozon_shortener/internal/api/url"
	"ozon_shortener/internal/config"
	"ozon_shortener/internal/repository"
	"ozon_shortener/internal/services/url"

	_ "github.com/lib/pq"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
)

const (
	successExitCode = 0
	failExitCode    = 1
)

// @title      URL shortener
// @version    1.0
// @description  URL shortener on go
// @host      localhost:8080
// @BasePath    /api/v1
func main() {
	os.Exit(run())
}

func run() (exitCode int) {
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

	dao := repository.NewRepository(db)

	urlSrv := url.New(dao, cfg.App.Url)

	urlAdapters := adapters.New(urlSrv)

	handler := handlers.New(urlAdapters)

	r := mux.NewRouter()

	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	r.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	apiRouter := r.PathPrefix("/api").Subrouter()

	urlRouter := apiRouter.PathPrefix("/url").Subrouter()

	urlRouter.HandleFunc("/generate", handler.CreateURL).Methods("POST")
	urlRouter.HandleFunc("/original", handler.GetOriginal).Methods("GET")

	slog.Info("start serve", slog.String("url", cfg.App.Url))
	if err := http.ListenAndServe(cfg.App.Url, r); err != nil {
		slog.Error("serve", "err", err)
		return failExitCode
	}

	// server := &http.Server{
	// 	Addr:    cfg.App.Url,
	// 	Handler: r,
	// }
	// lis, err := net.Listen("tcp", cfg.App.Url)
	// if err != nil {
	// 	slog.Error("take port", "err", err)
	// 	return failExitCode
	// }

	// slog.Info("start serve", slog.String("app url", cfg.App.Url))
	// if err := server.Serve(lis); err != nil {
	// 	slog.Error("take port", "err", err)
	// 	return failExitCode
	// }
	// return successExitCode

	return successExitCode
}

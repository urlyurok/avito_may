package main

import (
	"avitoTech/internal/app"
	"avitoTech/internal/config"
	"avitoTech/internal/repo"
	"avitoTech/internal/router"
	"avitoTech/internal/service"
	"avitoTech/internal/storage/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := app.SetupLogger(cfg.LogLevel)

	log.Info("config loaded", "log level", cfg.LogLevel)

	storage, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		log.Error("cannot init storage", err)
		os.Exit(1)
	}

	repositories := repo.NewRepos(storage)
	services := service.NewServices(repositories)

	log.Info("Initializing handlers and routes...")
	r := router.NewRouter(services)

	// Сообщаем, что сервер запускается
	log.Info("Server starting on port: " + cfg.HTTPServer.Adress)

	// Запускаем сервер и обрабатываем возможные ошибки
	if err := http.ListenAndServe(""+cfg.HTTPServer.Adress, r); err != nil {
		log.Error("Server failed to start", err)
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

}

package main

import (
	"BestMusicLibrary"
	"BestMusicLibrary/cfg"
	_ "BestMusicLibrary/docs"
	"BestMusicLibrary/internal/client"
	"BestMusicLibrary/internal/handler"
	"BestMusicLibrary/internal/repository"
	"BestMusicLibrary/internal/service"
	"BestMusicLibrary/migrations"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title MusicLibrary App
// @version 1.0
// @description API Server for MusicLibrary application
// @host localhost:8080
// @BasePath /
func main() {
	config := cfg.Get()
	db, err := repository.NewPostgresDb(repository.Config{Host: config.DbHost, Port: config.DbPort, UserName: config.DbUser, Password: config.DbPassword, DbName: config.DbName, SSLMode: config.DbSSLMode})
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(db)

	if err != nil {
		logrus.Fatal(err)
		return
	}

	dbMigrator := migrations.NewDbMigrator(db, "migrations")
	if err = dbMigrator.Migrate(); err != nil {
		logrus.Error(err)
		return
	}

	repos := repository.NewRepository(db)
	externalClient := client.NewExternalSongApiClient(config.ExternalApiClientUrl)
	mainService := service.NewService(repos, externalClient)
	hand := handler.NewHandler(mainService)
	srv := BestMusicLibrary.Server{}

	hand.InitRoutes()
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	http.Handle("/docs/", http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))))

	go func() {
		if err = srv.Run(config.ServerPort, nil); err != nil {
			logrus.Error(err)
		}
	}()

	logrus.Infof("listening on :%s", config.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Stop(ctx)
}

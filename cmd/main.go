package main

import (
	"BestMusicLibrary"
	"BestMusicLibrary/internal/cfg"
	"BestMusicLibrary/internal/client"
	"BestMusicLibrary/internal/handler"
	"BestMusicLibrary/internal/repository"
	"BestMusicLibrary/internal/service"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	config := cfg.Get()
	db, err := repository.NewPostgresDb(repository.Config{Host: config.DbHost, Port: config.DbPort, UserName: config.DbUser, Password: config.DbPassword, DbName: config.DbName, SSLMode: config.DbSSLMode})
	if err != nil {
		logrus.Fatal(err)
		return
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(db)

	repos := repository.NewRepository(db)
	externalClient := client.NewExternalSongApiClient(config.ExternalApiClientUrl)
	mainService := service.NewService(repos, externalClient)
	hand := handler.NewHandler(mainService)
	srv := BestMusicLibrary.Server{}

	hand.InitRoutes()

	go func() {
		if err = srv.Run(config.ServerPort, nil); err != nil {
			logrus.Fatal(err)
		}
	}()

	logrus.Infof("listening on :%s", config.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logrus.Info("shutting down server...")
	_ = srv.Stop(context.Background())
}

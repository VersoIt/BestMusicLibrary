package main

import (
	"BestMusicLibrary"
	"BestMusicLibrary/internal/handler"
	"BestMusicLibrary/internal/repository"
	"BestMusicLibrary/internal/service"
	"BestMusicLibrary/pkg/client"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	repos := repository.NewRepository()
	externalClient := client.NewExternalSongApiClient("https://external-api.com")
	serv := service.NewService(repos, externalClient)
	hand := handler.NewHandler(serv)
	hand.InitRoutes()
	srv := BestMusicLibrary.Server{}
	go func() {
		if err := srv.Run("8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Listening on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	if err := srv.Stop(context.Background()); err != nil {
		log.Println(err)
	}
}

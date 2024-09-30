package service

import (
	"BestMusicLibrary/internal/repository"
	"BestMusicLibrary/pkg/client"
	"BestMusicLibrary/pkg/model"
)

type Song interface {
	GetSongs(group, song string, page, limit int) ([]model.Song, error)
	GetSongText(id int64, page, limit int) (string, error)
	DeleteSong(id int64) error
	UpdateSong(song model.Song) error
	AddSong(song model.Song) error
}

type Service struct {
	Song   Song
	client client.ExternalSongApiClient
}

func NewService(repos *repository.Repository, externalClient client.ExternalSongApiClient) *Service {
	return &Service{Song: NewSongService(repos.Song), client: externalClient}
}

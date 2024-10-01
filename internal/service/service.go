package service

import (
	"BestMusicLibrary/internal/client"
	"BestMusicLibrary/internal/model"
	"BestMusicLibrary/internal/repository"
)

type Song interface {
	GetSongs(group, song string, page, limit int) ([]model.Song, error)
	GetSongVerses(id int64, page, limit int) ([]model.Verse, error)
	DeleteSong(id int64) error
	UpdateSong(song model.Song, text string) error
	AddSong(song model.Song) (int64, error)
}

type Service struct {
	Song Song
}

func NewService(repos *repository.Repository, externalClient client.ExternalSongApiClient) *Service {
	return &Service{Song: NewSongService(repos.Song, externalClient)}
}

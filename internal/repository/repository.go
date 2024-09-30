package repository

import (
	"BestMusicLibrary/pkg/model"
)

type Song interface {
	GetSongs(group, song string, page, limit int) ([]model.Song, error)
	GetSongText(page, limit int) (string, error)
	DeleteSong(id int64) error
	UpdateSong(song model.Song) error
	AddSong(song model.Song) error
}

type Repository struct {
	Song Song
}

func NewRepository() *Repository {
	return &Repository{}
}

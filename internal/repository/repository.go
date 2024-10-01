package repository

import (
	"BestMusicLibrary/internal/model"
	"github.com/jmoiron/sqlx"
)

type Song interface {
	GetSongs(group, song string, page, limit int) ([]model.Song, error)
	GetSongVerses(id int64, page, limit int) ([]model.Verse, error)
	DeleteSong(id int64) error
	UpdateSong(song model.Song) error
	AddSong(song model.Song) (int64, error)
}

type Repository struct {
	Song Song
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Song: &SongPostgresRepository{db: db}}
}

package repository

import (
	"BestMusicLibrary/pkg/model"
)

type SongPostgresRepository struct {
}

func (*SongPostgresRepository) NewSongRepository() *SongPostgresRepository {
	return &SongPostgresRepository{}
}

func (s *SongPostgresRepository) AddSong(song *model.Song) error {
	return nil
}

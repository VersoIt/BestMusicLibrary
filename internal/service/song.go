package service

import (
	"BestMusicLibrary/internal/repository"
	"BestMusicLibrary/pkg/client"
	"BestMusicLibrary/pkg/model"
	"log"
)

type SongService struct {
	songRepos  repository.Song
	songClient client.ExternalSongApiClient
}

func NewSongService(repos repository.Song) *SongService {
	return &SongService{songRepos: repos}
}

// GetSongs Получение данных библиотеки с фильтрацией по всем полям и пагинацией
func (s *SongService) GetSongs(group, song string, page, limit int) ([]model.Song, error) {
	return s.songRepos.GetSongs(group, song, page, limit)
}

// GetSongText Получение текста песни с пагинацией по куплетам
func (s *SongService) GetSongText(id int64, page, limit int) (string, error) {
	return s.songRepos.GetSongText(page, limit)
}

// DeleteSong Удаление песни
func (s *SongService) DeleteSong(id int64) error {
	return s.songRepos.DeleteSong(id)
}

// UpdateSong Изменение песни
func (s *SongService) UpdateSong(song model.Song) error {
	return s.songRepos.UpdateSong(song)
}

// AddSong Добавление песни
func (s *SongService) AddSong(song model.Song) error {
	enrichedSong, err := s.enrichSongWithAPI(song)
	if err != nil {
		return err
	}

	return s.songRepos.AddSong(enrichedSong)
}

// EnrichSongWithAPI Обогащение данных с использованием стороннего сервиса
func (s *SongService) enrichSongWithAPI(song model.Song) (enrichedSong model.Song, err error) {
	log.Printf("enriching songRepos data for group: %s, songRepos: %s\n", song.Group, song.Name)

	songDetails, err := s.songClient.FetchSongDetails(song.Group, song.Name)
	if err != nil {
		return song, err
	}

	enrichedSong.ReleaseDate = songDetails.ReleaseDate
	enrichedSong.Text = songDetails.Text
	enrichedSong.Link = songDetails.Link

	return enrichedSong, nil
}

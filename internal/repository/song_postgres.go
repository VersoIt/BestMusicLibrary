package repository

import (
	"BestMusicLibrary/internal/model"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SongPostgresRepository struct {
	db *sqlx.DB
}

func (*SongPostgresRepository) NewSongRepository() *SongPostgresRepository {
	return &SongPostgresRepository{}
}

func (s *SongPostgresRepository) GetSongs(group, songName string, page, limit int) ([]model.Song, error) {
	offset := page * limit
	rows, err := s.db.Query(`SELECT * FROM songs WHERE (LENGTH($1) > 0 AND group_name ILIKE '%' || $1 || '%') OR (LENGTH($2) > 0 AND song_title ILIKE '%' || $2 || '%') ORDER BY id LIMIT $3 OFFSET $4`, group, songName, limit, offset)
	if err != nil {
		return nil, err
	}

	songs := make([]model.Song, 0)
	for rows.Next() {
		var song model.Song
		err = rows.Scan(&song.Id, &song.Group, &song.Name, &song.ReleaseDate, &song.Link, &song.CreatedAt, &song.UpdatedAt)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (s *SongPostgresRepository) GetSongVerses(id int64, page, limit int) ([]model.Verse, error) {
	offset := page * limit
	rows, err := s.db.Query(`SELECT verse_number, text FROM verses WHERE song_id = $1 ORDER BY ID LIMIT $2 OFFSET $3`, id, limit, offset)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(rows)

	if err != nil {
		return nil, err
	}

	verses := make([]model.Verse, 0)
	for rows.Next() {
		var verse model.Verse
		err = rows.Scan(&verse.VerseNumber, &verse.Text)
		if err != nil {
			return nil, err
		}
		verses = append(verses, verse)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return verses, nil
}

func (s *SongPostgresRepository) DeleteSong(id int64) error {
	_, err := s.db.Exec(`DELETE FROM songs WHERE id = $1`, id)
	return err
}

func (s *SongPostgresRepository) UpdateSong(song model.Song) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`UPDATE songs SET group_name = $1, song_title = $2, release_date = $3, link = $4, created_at = $5, updated_at = NOW() WHERE id = $6`,
		song.Group, song.Name, song.ReleaseDate, song.Link, song.CreatedAt, song.Id)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = s.db.Exec(`DELETE FROM verses WHERE song_id = $1`, song.Id)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for index, verse := range song.Verses {
		_, err = s.db.Exec(`INSERT INTO verses(song_id, verse_number, text) VALUES($1, $2, $3)`, song.Id, index, verse.Text)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *SongPostgresRepository) AddSong(song model.Song) (int64, error) {
	var songId int64
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	err = s.db.QueryRow(`
		INSERT
		INTO
		songs(group_name, song_title, release_date, link)
		VALUES($1, $2, $3, $4) RETURNING
		id
		`,
		song.Group, song.Name, song.ReleaseDate, song.Link).Scan(&songId)

	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	for _, verse := range song.Verses {
		_, err = s.db.Exec(`
		INSERT
		INTO
		verses(song_id, verse_number, text)
		VALUES($1, $2, $3)`,
			songId, verse.VerseNumber, verse.Text)

		if err != nil {
			_ = tx.Rollback()
			return songId, err
		}
	}

	return songId, tx.Commit()
}

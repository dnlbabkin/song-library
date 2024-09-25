package repository

import (
	"database/sql"
	"fmt"
	"log"
	"song-library/internal/models"
	"strings"
)

type SongRepository struct {
	DB *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	log.Println("[INFO] Initializing SongRepository")
	return &SongRepository{DB: db}
}

func (r *SongRepository) AddSong(song models.Song) error {
	log.Printf("[INFO] Adding song: %s - %s\n", song.GroupName, song.SongName)
	query := `INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)

	if err != nil {
		log.Printf("[ERROR] Failed to add song: %v\n", err)
		return err
	}

	log.Printf("[INFO] Successfully added song: %s - %s\n", song.GroupName, song.SongName)

	return nil
}

func (r *SongRepository) GetSong(id int) (models.Song, error) {
	log.Printf("[DEBUG] Fetching song with ID: %d\n", id)
	var song models.Song
	query := `SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE id=$1`
	err := r.DB.QueryRow(query, id).Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)

	if err != nil {
		log.Printf("[ERROR] Failed to get song with ID %d: %v\n", id, err)
		return song, err
	}

	log.Printf("[INFO] Successfully fetched song: %s - %s\n", song.GroupName, song.SongName)

	return song, nil
}

func (r *SongRepository) DeleteSong(id int) error {
	log.Printf("[INFO] Deleting song with ID: %d\n", id)
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.DB.Exec(query, id)

	if err != nil {
		log.Printf("[ERROR] Failed to delete song with ID %d: %v\n", id, err)
		return err
	}

	log.Printf("[INFO] Successfully deleted song with ID: %d\n", id)

	return nil
}

func (r *SongRepository) UpdateSong(id int, song models.Song) error {
	log.Printf("[INFO] Updating song with ID: %d\n", id)
	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6`
	_, err := r.DB.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, id)

	if err != nil {
		log.Printf("[ERROR] Failed to update song with ID %d: %v\n", id, err)
		return err
	}

	log.Printf("[INFO] Successfully updated song with ID: %d\n", id)

	return nil
}

func (r *SongRepository) GetSongsFiltered(group, song string, limit, offset int) ([]models.Song, error) {
	log.Printf("[DEBUG] Fetching songs with filter - group: %s, song: %s, limit: %d, offset: %d\n", group, song, limit, offset)
	var songs []models.Song
	var filters []string
	var params []interface{}

	if group != "" {
		filters = append(filters, "group_name ILIKE $"+fmt.Sprint(len(params)+1))
		params = append(params, "%"+group+"%")
	}
	if song != "" {
		filters = append(filters, "song_name ILIKE $"+fmt.Sprint(len(params)+1))
		params = append(params, "%"+song+"%")
	}

	query := "SELECT id, group_name, song_name, release_date, text, link FROM songs"
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(params)+1, len(params)+2)
	params = append(params, limit, offset)

	rows, err := r.DB.Query(query, params...)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch songs: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			log.Printf("[ERROR] Error scanning song: %v\n", err)
			continue
		}
		songs = append(songs, song)
	}

	log.Printf("[INFO] Successfully fetched %d songs\n", len(songs))

	return songs, nil
}

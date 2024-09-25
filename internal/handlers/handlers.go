package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"song-library/internal/models"
	"song-library/internal/repository"
	"song-library/internal/service"
	"strconv"
)

type SongHandler struct {
	Repo repository.SongRepository
	Svc  service.ExternalAPIService
}

// AddSong godoc
// @Summary Add a new song
// @Description Add a new song with details fetched from external API
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Song info"
// @Success 200 {object} models.Song
// @Failure 400 {object} gin.H
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var input models.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("[INFO] Fetching external API data for group: %s, song: %s\n", input.GroupName, input.SongName)
	songDetail, err := h.Svc.FetchSongDetails(input.GroupName, input.SongName)
	if err != nil {
		log.Printf("[ERROR] Failed to get song detail from external API: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song data"})
		return
	}

	song := models.Song{
		GroupName:   input.GroupName,
		SongName:    input.SongName,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	if err := h.Repo.AddSong(song); err != nil {
		log.Printf("[ERROR] Failed to add song to the database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add song"})
		return
	}

	log.Printf("[INFO] Successfully added song: %s - %s\n", song.GroupName, song.SongName)
	c.JSON(http.StatusCreated, song)
}

// GetSong @Summary Get song by ID
// @Description Retrieve a song by its ID from the database.
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.Song
// @Failure 400 {object} gin.H {"error": "Invalid song ID"}
// @Failure 404 {object} gin.H {"error": "Song not found"}
// @Router /songs/{id} [get]
func (h *SongHandler) GetSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Invalid song ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	log.Printf("[INFO] Fetching song with ID: %d\n", id)
	song, err := h.Repo.GetSong(id)
	if err != nil {
		log.Printf("[ERROR] Failed to get song with ID %d: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	log.Printf("[INFO] Successfully fetched song with ID: %d\n", id)
	c.JSON(http.StatusOK, song)
}

// DeleteSong @Summary Delete song by ID
// @Description Delete a song by its ID from the database.
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 204
// @Failure 400 {object} gin.H {"error": "Invalid song ID"}
// @Failure 500 {object} gin.H {"error": "Failed to delete song"}
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Invalid song ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	log.Printf("[INFO] Deleting song with ID: %d\n", id)
	if err := h.Repo.DeleteSong(id); err != nil {
		log.Printf("[ERROR] Failed to delete song with ID %d: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	log.Printf("[INFO] Successfully deleted song with ID: %d\n", id)
	c.Status(http.StatusNoContent)
}

// UpdateSong @Summary Update song by ID
// @Description Update song details by its ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song details"
// @Success 200
// @Failure 400 {object} gin.H {"error": "Invalid song ID or input"}
// @Failure 500 {object} gin.H {"error": "Failed to update song"}
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[ERROR] Invalid song ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var input models.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("[INFO] Updating song with ID: %d\n", id)
	if err := h.Repo.UpdateSong(id, input); err != nil {
		log.Printf("[ERROR] Failed to update song with ID %d: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	log.Printf("[INFO] Successfully updated song with ID: %d\n", id)
	c.Status(http.StatusOK)
}

// GetSongsFiltered @Summary Get songs with filtering and pagination
// @Description Retrieve songs with optional filtering by group and song name, with pagination support.
// @Tags songs
// @Produce json
// @Param group query string false "Group name"
// @Param song query string false "Song name"
// @Param limit query int false "Number of songs to return" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.Song
// @Failure 500 {object} gin.H {"error": "Failed to fetch songs"}
// @Router /songs [get]
func (h *SongHandler) GetSongsFiltered(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	log.Printf("[DEBUG] Fetching songs with filters - group: %s, song: %s, limit: %d, offset: %d\n", group, song, limit, offset)

	songs, err := h.Repo.GetSongsFiltered(group, song, limit, offset)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch songs: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch songs"})
		return
	}

	log.Printf("[INFO] Successfully fetched %d songs\n", len(songs))
	c.JSON(http.StatusOK, songs)
}

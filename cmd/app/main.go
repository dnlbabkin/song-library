package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"song-library/internal/config"
	"song-library/internal/handlers"
	"song-library/internal/repository"
	"song-library/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres",
		"host="+cfg.DBHost+
			" port="+cfg.DBPort+
			" user="+cfg.DBUser+
			" password="+cfg.DBPassword+
			" dbname="+cfg.DBName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewSongRepository(db)
	svc := service.NewExternalAPIService(cfg.ExternalAPI)
	handler := handlers.SongHandler{Repo: *repo, Svc: *svc}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/songs", handler.AddSong)
	r.GET("/songs/:id", handler.GetSong)
	r.DELETE("/songs/:id", handler.DeleteSong)
	r.PUT("/songs/:id", handler.UpdateSong)
	r.GET("/songs", handler.GetSongsFiltered)

	r.Run(":" + cfg.AppPort)
}

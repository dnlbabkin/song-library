package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExternalAPIService struct {
	APIUrl string
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewExternalAPIService(apiUrl string) *ExternalAPIService {
	return &ExternalAPIService{APIUrl: apiUrl}
}

func (s *ExternalAPIService) FetchSongDetails(group, song string) (*SongDetail, error) {
	resp, err := http.Get(fmt.Sprintf("%s?group=%s&song=%s", s.APIUrl, group, song))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch song details, status code: %d", resp.StatusCode)
	}

	var songDetail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, err
	}

	return &songDetail, nil
}

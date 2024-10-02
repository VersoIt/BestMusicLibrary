package client

import (
	"BestMusicLibrary/internal/service"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type songApiResponse struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewExternalSongApiClient(baseUrl string) *ExternalSongApiClient {
	return &ExternalSongApiClient{baseUrl: baseUrl}
}

type ExternalSongApiClient struct {
	baseUrl string
}

func (c *ExternalSongApiClient) FetchSongDetails(group, song string) (service.SongFetchData, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.baseUrl, group, song)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return service.SongFetchData{}, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logrus.Errorf("error closing client response body from %s\n", url)
		}
	}(resp.Body)

	var songResponse songApiResponse
	if err = json.NewDecoder(resp.Body).Decode(&songResponse); err != nil {
		return service.SongFetchData{}, err
	}

	return service.SongFetchData{ReleaseDate: songResponse.ReleaseDate, Link: songResponse.Link, Text: songResponse.Text}, nil
}

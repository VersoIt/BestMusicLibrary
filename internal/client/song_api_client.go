package client

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type SongApiResponse struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewExternalSongApiClient(baseUrl string) ExternalSongApiClient {
	return &ExternalSongApiClientImpl{baseUrl: baseUrl}
}

type ExternalSongApiClient interface {
	FetchSongDetails(group, song string) (*SongApiResponse, error)
}

type ExternalSongApiClientImpl struct {
	baseUrl string
}

func (c *ExternalSongApiClientImpl) FetchSongDetails(group, song string) (*SongApiResponse, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.baseUrl, group, song)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logrus.Errorf("error closing client response body from %s\n", url)
		}
	}(resp.Body)

	var songResponse SongApiResponse
	if err = json.NewDecoder(resp.Body).Decode(&songResponse); err != nil {
		return nil, err
	}

	return &songResponse, nil
}

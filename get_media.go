package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/grafov/m3u8"
)

type Media struct {
	ID  int
	URI string
}

func (app *App) GetMedia(id, resolution, uri string) ([]*Media, error) {
	client := app.Client
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Add("host", "surrit.com")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	dat, list, err := m3u8.DecodeFrom(res.Body, true)
	if err != nil {
		return nil, err
	}

	switch list {
	case m3u8.MEDIA:
		mediapl := dat.(*m3u8.MediaPlaylist)
		uris := []*Media{}

		for idx, segment := range mediapl.Segments {
			if segment == nil {
				continue
			}

			u := segment.URI
			resUrl := fmt.Sprintf("https://surrit.com/%s/%s/%s", id, resolution, u)
			med := Media{URI: resUrl, ID: idx}
			uris = append(uris, &med)
		}

		return uris, nil
	}

	return nil, errors.New("media not found")
}

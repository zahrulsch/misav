package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/grafov/m3u8"
)

type URI = string
type Res = string

func getPlaylist(id string) (URI, Res, error) {
	uriLike := fmt.Sprintf("https://surrit.com/%v/playlist.m3u8", id)
	transport := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", uriLike, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Add("host", "https://missav.com")
	req.Header.Add("referer", "https://missav.com/dm5/en")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	defer res.Body.Close()

	dat, list, err := m3u8.DecodeFrom(res.Body, true)
	if err != nil {
		return "", "", err
	}

	switch list {
	case m3u8.MASTER:
		masterpl := dat.(*m3u8.MasterPlaylist)
		lenVars := len(masterpl.Variants)

		if lenVars < 1 {
			return "", "", errors.New("empty variants")
		}

		// max := masterpl.Variants[0]
		max := masterpl.Variants[lenVars-1]
		res := ""

		maxRes := strings.Split(max.URI, "/")

		if len(maxRes) > 0 {
			res = maxRes[0]
		}

		uri := fmt.Sprintf("https://surrit.com/%s/%s", id, max.URI)
		return uri, res, nil
	}

	return "", "", errors.New("variants not found")

}

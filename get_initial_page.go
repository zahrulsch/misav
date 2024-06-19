package main

import (
	"errors"
	"io"
	"net/http"
	"regexp"
)

type VideoID = string

func (app *App) GetInitialPage(url string) (VideoID, error) {
	client := app.Client

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("host", "https://missav.com")
	req.Header.Add("referer", "https://missav.com/dm5/en")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Add("upgrade-insecure-requests", "1")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`https:\\\/\\\/\w+\.com\\\/([\w-]+)\\\/seek\\\/_`)
	match := re.FindStringSubmatch(string(responseData))

	if len(match) < 2 {
		return "", errors.New("id not found")
	}

	if match[1] == "" {
		return "", errors.New("id not found")
	}

	return match[1], nil
}

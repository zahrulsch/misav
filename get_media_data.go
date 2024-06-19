package main

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
)

var transport = &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13}, ForceAttemptHTTP2: false}
var client = &http.Client{Transport: transport}

func getMediaData(file *os.File, uri string) error {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Add("host", "surrit.com")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if _, err := io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}

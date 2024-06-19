package main

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	MaxConcurrentMedia int    `toml:"max_concurrent_media,omitempty,commented" json:"max_concurrent_media"`
	OutDir             string `toml:"out_dir,omitempty,commented" json:"out_dir"`
	TargetFile         string `toml:"target_file,omitempty,commented" json:"target_file"`
}

func (cfg *Config) JSONString() string {
	s, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func LoadConfig() *Config {
	p, _ := filepath.Abs("config.toml")
	file, err := os.ReadFile(p)

	if err != nil {
		panic(err)
	}

	cfg := Config{}

	if err := toml.Unmarshal(file, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

type App struct {
	Config *Config
	Client *http.Client
}

func NewApp() *App {
	var transport = &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13}, ForceAttemptHTTP2: false}
	var client = &http.Client{Transport: transport}

	app := App{
		Config: LoadConfig(),
		Client: client,
	}

	return &app
}

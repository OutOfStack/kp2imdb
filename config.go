package main

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	localeEn = "en"
	localeRu = "ru"
)

type config struct {
	Cookie     string `json:"imdb_cookie"`
	OMDbAPIKey string `json:"omdb_api_key"`
	Locale     string `json:"locale"`
}

// loadConfig reads and validates configuration from settings.json
func loadConfig() (*config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.OMDbAPIKey == "" {
		return nil, errors.New("OMDb API Key not set")
	}
	if cfg.Cookie == "" {
		return nil, errors.New("IMDb Cookie not set")
	}
	if cfg.Locale != localeEn && cfg.Locale != localeRu {
		return nil, errors.New("invalid locale: should be either 'en' or 'ru'")
	}

	return &cfg, nil
}

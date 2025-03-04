package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const omdbURL = "https://www.omdbapi.com/"

const omdbAPIMovieNotFoundError = "Movie not found!"

type omdbResponse struct {
	Response string `json:"Response"`
	IMDbID   string `json:"imdbID"`
	Title    string `json:"title"`
	Error    string `json:"Error"`
}

type omdbMovie struct {
	IMDbID string
	Title  string
}

// omdb client
type omdb struct {
	apiKey string
}

// searchMovie uses the OMDb API to find the IMDb ID using original name and release year
func (c *omdb) searchMovie(originalName string, year int) (omdbMovie, error) {
	query := url.Values{}
	query.Set("apikey", c.apiKey)
	query.Set("t", originalName)
	if year != 0 {
		query.Set("y", strconv.Itoa(year))
	}
	fullURL := omdbURL + "?" + query.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return omdbMovie{}, fmt.Errorf("OMDb API request error: %v", err)
	}
	defer resp.Body.Close()

	var omdbResp omdbResponse
	err = json.NewDecoder(resp.Body).Decode(&omdbResp)
	if err != nil {
		return omdbMovie{}, fmt.Errorf("decode OMDb API response: %v", err)
	}

	if omdbResp.Response != "True" {
		if omdbResp.Error == omdbAPIMovieNotFoundError {
			return omdbMovie{}, errNotFound
		}
		return omdbMovie{}, fmt.Errorf("OMDb API error: %s", omdbResp.Error)
	}

	return omdbMovie{
		IMDbID: omdbResp.IMDbID,
		Title:  omdbResp.Title,
	}, nil
}

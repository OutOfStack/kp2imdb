package main

type movie struct {
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Rating       string `json:"my_rating"`
	Year         string `json:"year"`
}

type rateError struct {
	movie
	Error string `json:"error"`
	Date  string `json:"date"`
}

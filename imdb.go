package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	imdbAPIURL = "https://api.graphql.imdb.com/"
)

type graphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// imdb client
type imdb struct {
	Client *http.Client
	Cookie string
}

// updateRating sends a GraphQL mutation to update the title rating
func (i *imdb) updateRating(movieID string, rating int) error {
	reqBody := graphQLRequest{
		Query: `
mutation UpdateTitleRating($rating: Int!, $titleId: ID!) {
  rateTitle(input: {rating: $rating, titleId: $titleId}) {
    rating {
      value
    }
  }
}`,
		OperationName: "UpdateTitleRating",
		Variables: map[string]interface{}{
			"rating":  rating,
			"titleId": movieID,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshall GraphQL request: %v", err)
	}

	req, err := http.NewRequest("POST", imdbAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create POST request for %s: %v", movieID, err)
	}

	req.Header.Set("accept", "application/graphql+json, application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", i.Cookie)

	resp, err := i.Client.Do(req)
	if err != nil {
		log.Printf("Error sending POST request for movie %s: %v", movieID, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, rErr := io.ReadAll(resp.Body)
		if rErr != nil {
			return fmt.Errorf("read response body: %v", rErr)
		}
		log.Printf("Received status code: %d with body: %s for %s when updating rating", resp.StatusCode, body, movieID)
		return fmt.Errorf("received status code: %d with body: %s for %s", resp.StatusCode, body, movieID)
	}

	return nil
}

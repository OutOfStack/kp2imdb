package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Reading config: %v", err)
	}

	// read movies data from data.json
	file, err := os.Open("data.json")
	if err != nil {
		log.Fatalf("Error opening data.json: %v", err)
	}
	defer file.Close()

	var movies []movie
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&movies); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	p := processor{
		omdbClient: &omdb{
			apiKey: cfg.OMDbAPIKey,
		},
		imdbClient: &imdb{
			Client: &http.Client{},
			Cookie: cfg.Cookie,
		},
	}

	// process each movie entry in reverse order as movie list collected in reversed chronological order
	for i := len(movies) - 1; i >= 0; i-- {
		err = p.processMovie(&movies[i])
		if err != nil {
			rateEntry := rateError{
				movie: movies[i],
				Date:  time.Now().Format(time.RFC3339),
				Error: getErrorMsg(err, cfg.Locale),
			}

			err = appendErrorToFile(rateEntry)
			if err != nil {
				log.Printf("Error appending to warnings.json: %v", err)
			}
		}
	}
}

// appendErrorToFile reads the existing warnings.json file, appends the new error entry,
// and writes back the updated array without replacing the previous entries
func appendErrorToFile(e rateError) error {
	var errorsSlice []rateError

	// read the existing warnings.json file
	data, err := os.ReadFile("warnings.json")
	if err == nil {
		err = json.Unmarshal(data, &errorsSlice)
		if err != nil {
			return fmt.Errorf("unmarshall warnings.json: %v", err)
		}
	} else {
		// if the file doesn't exist, initialize an empty slice
		if os.IsNotExist(err) {
			errorsSlice = []rateError{}
		} else {
			return fmt.Errorf("read warnings.json: %v", err)
		}
	}

	errorsSlice = append(errorsSlice, e)

	updatedData, err := json.MarshalIndent(errorsSlice, "", "	")
	if err != nil {
		return fmt.Errorf("marshall errors slice: %v", err)
	}

	// write the updated data back to warnings.json
	err = os.WriteFile("warnings.json", updatedData, 0600)
	if err != nil {
		return fmt.Errorf("write warnings.json: %v", err)
	}

	return nil
}

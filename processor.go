package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xrash/smetrics"
)

const (
	similarityThreshold = 0.8
)

var (
	errNoRating     = errors.New("no rating")
	errNotFound     = errors.New("not found")
	errUpdateRating = errors.New("update rating")
	errOMDbApiLimit = errors.New("OMDb API limit")
)

type maybeIncorrectMovieError struct {
	title     string
	year      int
	searchAs  string
	imdbTitle string
	imdbID    string
}

func (e *maybeIncorrectMovieError) Error() string {
	return "maybe incorrect movie"
}

type processor struct {
	omdbClient *omdb
	imdbClient *imdb
}

func (p *processor) processMovie(m *movie) error {
	// skip entries without name and original name
	if m.Name == "" && m.OriginalName == "" {
		return nil
	}

	rating, _ := strconv.Atoi(m.Rating)
	if rating <= 0 || rating > 10 {
		return errNoRating
	}

	originalName := m.OriginalName
	if m.OriginalName == "" {
		originalName = m.Name
	}

	var year int
	if m.Year != "" {
		year, _ = strconv.Atoi(m.Year)
	}

	// search for the movie on omdb
	movieInfo, sErr := p.searchMovie(originalName, m.Name, year)
	var miv *maybeIncorrectMovieError
	if sErr != nil && !errors.As(sErr, &miv) {
		fmt.Printf("❌ '%s' (%s) not found\n", originalName, m.Year)
		return sErr
	}

	// set movie rating
	err := p.imdbClient.updateRating(movieInfo.IMDbID, rating)
	if err != nil {
		log.Printf("Error updating rating for %s: %v", originalName, err)
		return errUpdateRating
	}

	fmt.Printf("✅ %s (%s), Rating: %s, IMDb Title: %s, IMDb link: https://www.imdb.com/title/%s\n", originalName, m.Year, m.Rating, movieInfo.Title, movieInfo.IMDbID)

	// return maybeIncorrectMovie in order to write potential incorrect movie into errors list
	return sErr
}

// searches for movie on omdb applying different alternating some movie data
func (p *processor) searchMovie(originalName, name string, year int) (omdbMovie, error) {
	// try to search by original name
	movieInfo, err := p.omdbClient.searchMovie(originalName, year)
	if err == nil && (movieInfo.Title == originalName) {
		return movieInfo, nil
	}
	if err != nil && !errors.Is(err, errNotFound) {
		log.Printf("Search '%s': %v", originalName, err)
		return omdbMovie{}, errOMDbApiLimit
	}

	// try to search by localized name
	searchAs := name
	movieInfo, err = p.omdbClient.searchMovie(searchAs, year)
	if err == nil && (movieInfo.Title == searchAs) {
		return movieInfo, nil
	}
	if err != nil && !errors.Is(err, errNotFound) {
		log.Printf("Search '%s': %v", searchAs, err)
		return omdbMovie{}, errOMDbApiLimit
	}

	// try to search by next year (often release year on imdb is 1 year later than on kp)
	if year != 0 {
		searchYear := year + 1
		movieInfo, err = p.omdbClient.searchMovie(originalName, searchYear)
		if err == nil && (movieInfo.Title == originalName || movieInfo.Title == name) {
			fmt.Printf("🤔 '%s' (%d) found on IMDb as '%s' (%d), link - https://www.imdb.com/title/%s. Check if it's correct\n",
				originalName, year, movieInfo.Title, searchYear, movieInfo.IMDbID)
			return movieInfo, &maybeIncorrectMovieError{title: originalName, year: year, searchAs: originalName, imdbTitle: movieInfo.Title, imdbID: movieInfo.IMDbID}
		}
		if err != nil && !errors.Is(err, errNotFound) {
			log.Printf("Search '%s': %v", originalName, err)
			return omdbMovie{}, errOMDbApiLimit
		}
	}

	// try to find by transliterated name. If not found, return original error
	searchAs = transliterateRussian(name)
	movieInfo, err = p.omdbClient.searchMovie(searchAs, year)
	if err != nil && !errors.Is(err, errNotFound) {
		log.Printf("Search '%s': %v", searchAs, err)
		return omdbMovie{}, errOMDbApiLimit
	}
	if errors.Is(err, errNotFound) {
		return omdbMovie{}, err
	}

	// decide if searchAs/original title are similar to result title.
	// there are some false positives in omdb results
	similarityX := smetrics.Jaro(searchAs, movieInfo.Title)
	similarityY := smetrics.Jaro(originalName, movieInfo.Title)
	if similarityX < similarityThreshold && similarityY < similarityThreshold {
		fmt.Printf("❗'%s' (%d) searched as '%s', found as IMDb title - '%s', link - https://www.imdb.com/title/%s. Not rated because potentially incorrect\n",
			originalName, year, searchAs, movieInfo.Title, movieInfo.IMDbID)
		return omdbMovie{}, errNotFound
	}
	log.Printf("🤔'%s' (%d) searched as '%s', found IMDb title - '%s'. Check if it's correct\n", originalName, year, searchAs, movieInfo.Title)

	return movieInfo, &maybeIncorrectMovieError{title: originalName, year: year, searchAs: searchAs, imdbTitle: movieInfo.Title, imdbID: movieInfo.IMDbID}
}

func getErrorMsg(err error, locale string) string {
	var mivErr *maybeIncorrectMovieError
	switch {
	case errors.Is(err, errNoRating) && locale == localeEn:
		return "Movie has no rating. Try to add it to Check-in list manually"
	case errors.Is(err, errNoRating) && locale == localeRu:
		return "У фильма нет оценки. Добавь в список Check-in вручную"
	case errors.Is(err, errNotFound) && locale == localeEn:
		return "Movie not found. Find and rate it manually"
	case errors.Is(err, errNotFound) && locale == localeRu:
		return "Фильм не найден. Оцени его вручную"
	case errors.Is(err, errUpdateRating) && locale == localeEn:
		return "IMDb API error. Try to update session. If error persist, you're welcome to create an issue on github"
	case errors.Is(err, errUpdateRating) && locale == localeRu:
		return "Ошибка IMDb API. Попробуй обновить куки IMDb в конфиге. Если не поможет, можешь создать issue на гитхабе"
	case errors.Is(err, errOMDbApiLimit) && locale == localeEn:
		return "OMDb API limit exceeded. Try again tomorrow"
	case errors.Is(err, errOMDbApiLimit) && locale == localeRu:
		return "Превышен лимит запросов к OMDb. Попробуй повторить завтра"
	case errors.As(err, &mivErr) && locale == localeEn:
		return fmt.Sprintf("Rating is added to IMDb, but check if movie is correct. Kinopoisk title - '%s' (%d), searched as '%s', IMDb title - '%s', Link - https://www.imdb.com/title/%s",
			mivErr.title, mivErr.year, mivErr.searchAs, mivErr.imdbTitle, mivErr.imdbID)
	case errors.As(err, &mivErr) && locale == localeRu:
		return fmt.Sprintf("Рейтинг проставлен на IMDb, но нужно убедиться, что оценен корректный фильм. Название на КП - '%s' (%d), поисковый запрос '%s', фильм на IMDb - '%s', ссылка - https://www.imdb.com/title/%s",
			mivErr.title, mivErr.year, mivErr.searchAs, mivErr.imdbTitle, mivErr.imdbID)
	default:
		return err.Error()
	}
}

func transliterateRussian(input string) string {
	// transliteration mapping for Russian characters
	transMap := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
		'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i",
		'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
		'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
		'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch",
		'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "",
		'э': "e", 'ю': "yu", 'я': "ya",

		// Uppercase handling
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D",
		'Е': "E", 'Ё': "Yo", 'Ж': "Zh", 'З': "Z", 'И': "I",
		'Й': "Y", 'К': "K", 'Л': "L", 'М': "M", 'Н': "N",
		'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T",
		'У': "U", 'Ф': "F", 'Х': "Kh", 'Ц': "Ts", 'Ч': "Ch",
		'Ш': "Sh", 'Щ': "Shch", 'Ъ': "", 'Ы': "Y", 'Ь': "",
		'Э': "E", 'Ю': "Yu", 'Я': "Ya",
	}

	var result strings.Builder
	for _, char := range input {
		// check if the character is a Russian letter
		if transliterated, exists := transMap[char]; exists {
			result.WriteString(transliterated)
		} else {
			// if not a Russian letter, keep the original character
			result.WriteRune(char)
		}
	}

	return result.String()
}

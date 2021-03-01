package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	baseURL         = "https://xkcd.com"
	jsonPath        = "info.0.json"
	dataDir         = "data/"
	defaultComicNum = 2423
)

type Comic struct {
	Num        int
	Title      string
	Transcript string
}

func main() {
	downloadFlag := flag.Bool("d", false, "Download comics data")

	flag.Parse()

	if *downloadFlag {
		downloadAllComics()
	}
	index := buildSearchIndex()

	term := os.Args[len(os.Args)-1]
	search(term, index)
}

func downloadAllComics() {
	body, err := requestComic(comicUrl(""))
	var maxComicNum int
	if err != nil {
		fmt.Printf("Could not get most recent comic. Defaulting to #%d\n", defaultComicNum)
		maxComicNum = defaultComicNum
	} else {
		var comic Comic
		if err := json.NewDecoder(body).Decode(&comic); err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}
		maxComicNum = comic.Num
	}

	for i := maxComicNum; i > 0; i-- {
		comicNum := strconv.Itoa(i)
		saveLoc := dataDir + comicNum + ".json"

		// Download a comic only if it is not already on disk
		if _, err := os.Stat(saveLoc); os.IsNotExist(err) {
			fmt.Printf("Downloading xkcd #%s\n", comicNum)
			downloadErr := downloadComic(comicUrl(comicNum), saveLoc)
			fmt.Printf("%v", downloadErr)
		}
	}
}

func comicUrl(comicNum string) string {
	return strings.Join([]string{baseURL, comicNum, jsonPath}, "/")
}

func downloadComic(url string, saveLoc string) error {
	data, reqErr := requestComic(url)
	if reqErr != nil {
		return reqErr
	}

	saveErr := saveComic(saveLoc, data)
	return saveErr
}

func requestComic(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to download %s: %s\n", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Failed to download %s: %s\n", url, resp.Status)
	}

	return resp.Body, nil
}

func saveComic(location string, data io.ReadCloser) error {
	out, createErr := os.Create(location)
	if createErr != nil {
		return createErr
	}

	_, copyErr := io.Copy(out, data)
	return copyErr
}

func buildSearchIndex() map[string]map[int]bool {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	index := make(map[string]map[int]bool)
	for _, file := range files {
		comic := loadComicFromFile(file.Name())
		searchableText := comic.Title + comic.Transcript
		cleanedText := cleanText(searchableText)
		for _, word := range strings.Split(cleanedText, " ") {
			word = strings.TrimSpace(word)
			if index[word] == nil {
				index[word] = make(map[int]bool)
			}
			index[word][comic.Num] = true
		}
	}

	return index
}

func loadComicFromFile(fileName string) *Comic {
	data, err := ioutil.ReadFile(dataDir + fileName)
	if err != nil {
		log.Fatalf("Failed to open %v", err)
	}

	var comic Comic
	if err := json.Unmarshal(data, &comic); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return &comic
}

func cleanText(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "\n", " ")

	// There's always a "Title text:" description line
	text = strings.ReplaceAll(text, "title text:", "")

	// Remove all non-alphanumeric characters
	re := regexp.MustCompile(`[^a-zA-Z\d\s]`)
	text = re.ReplaceAllLiteralString(text, " ")

	return text
}

func search(term string, index map[string]map[int]bool) {
	cleanedTerm := cleanText(term)
	results, found := index[cleanedTerm]
	if !found {
		fmt.Printf("Search term: '%s' not found.\n", term)
		return
	}

	resultQuantifier := "result"
	numResults := len(results)
	if numResults != 1 {
		resultQuantifier += "s"
	}
	fmt.Printf("%d %s for '%s'\n", len(results), resultQuantifier, term)
	for num, _ := range results {
		printSearchResult(num)
	}
}

func printSearchResult(num int) {
	comicNum := strconv.Itoa(num)
	comic := loadComicFromFile(comicNum + ".json")
	url := comicUrl(comicNum)
	padding := fmt.Sprintf("%*s", len(url), "=")
	fmt.Printf("\n%s\n%s\n%s\n", url, strings.ReplaceAll(padding, " ", "="), comic.Transcript)
}

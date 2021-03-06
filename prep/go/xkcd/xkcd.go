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
	defaultDataDir  = "data/"
	defaultComicNum = 2430
)

type Comic struct {
	Num        int
	Title      string
	Transcript string
}
type resultSet map[int]bool
type searchIndex map[string]resultSet

var dataDir = defaultDataDir

// Expose as a global to facilitate capturing output for testing, as
// demonstrated in "The Go Programming Language" 11.2.2 (pg. 309)
var out io.Writer = os.Stdout

// Expose the requestComic function as a global to stub out the HTTP request
// functionality. Again, this technique was taken from "The Go Programming
// Language" 11.2.3 (pg. 312)
var requestComic = func(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error downloading %s: %s\n", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("Error downloading %s: %s\n", url, resp.Status)
	}

	return resp.Body, nil
}

func main() {
	downloadFlag := flag.Bool("d", false, "Download missing comics")
	saveDirFlag := flag.String("l", defaultDataDir, "Downloaded comics location")

	flag.Parse()

	if *saveDirFlag != defaultDataDir {
		if _, err := ioutil.ReadDir(*saveDirFlag); os.IsNotExist(err) {
			fmt.Fprintf(out, "Directory %s does not exist, exiting.\n", *saveDirFlag)
			return
		}
		dataDir = *saveDirFlag

		// Append the trailing forward slash if missing
		if dataDir[len(dataDir)-1] != '/' {
			dataDir += "/"
		}
	}

	if *downloadFlag {
		downloadAllComics(getMaxComicNum())
	}

	index := buildSearchIndex()

	// Grab the first non-flag argument
	term := flag.Arg(0)
	if term == "" {
		fmt.Fprintf(out, "No search term provided.\n")
		flag.PrintDefaults()
		return
	}

	results := search(term, index)
	printSearchResults(results, term)
}

func downloadAllComics(maxComicNum int) {
	for i := maxComicNum; i > 0; i-- {
		// This is an Easter Egg -- there is no Comic #404
		if i == 404 {
			continue
		}
		comicNum := strconv.Itoa(i)
		saveLoc := dataDir + comicNum + ".json"

		// Download a comic only if it is not already on disk
		if _, err := os.Stat(saveLoc); os.IsNotExist(err) {
			if err = downloadComic(comicUrl(comicNum), saveLoc); err != nil {
				fmt.Fprintln(out, err.Error())
			}
		}
	}
}

func buildSearchIndex() searchIndex {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	index := make(searchIndex)
	for _, file := range files {
		comic := loadComicFromFile(file.Name())
		searchableText := comic.Title + "\n" + comic.Transcript
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

func search(term string, index searchIndex) resultSet {
	cleanedTerm := cleanText(term)
	results, found := index[cleanedTerm]
	if !found {
		fmt.Fprintf(out, "Search term: '%s' not found.\n", term)
	}

	return results
}

func getMaxComicNum() int {
	// The most recent comic is at: xkcd.com/info.0.json, so a blank string will get it for us
	body, err := requestComic(comicUrl(""))
	if err != nil {
		fmt.Fprintf(out, "Could not get most recent comic. Defaulting to #%d\n", defaultComicNum)
		return defaultComicNum
	}

	var comic Comic
	if err = json.NewDecoder(body).Decode(&comic); err != nil {
		fmt.Fprintf(out, "Error parsing JSON for most recent comic. Defaulting to #%d\n", defaultComicNum)
		return defaultComicNum
	}
	return comic.Num
}

func comicUrl(comicNum string) string {
	return strings.Join([]string{baseURL, comicNum, jsonPath}, "/")
}

func downloadComic(url string, saveLoc string) error {
	fmt.Fprintf(out, "Downloading %s\n", url)
	data, err := requestComic(url)
	if err != nil {
		return err
	}

	err = saveComic(saveLoc, data)
	if err != nil {
		return fmt.Errorf("Error saving comic to %s: %s\n", saveLoc, err)
	}
	return nil
}

func saveComic(location string, data io.ReadCloser) error {
	out, err := os.Create(location)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, data)
	return err
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

	// There's always a "Title text:" description line
	text = strings.ReplaceAll(text, "title text:", "")

	// Remove all non-alphanumeric characters
	re := regexp.MustCompile(`[^a-zA-Z\d\s]`)
	text = re.ReplaceAllLiteralString(text, "")

	// Replace mutliple whitespace with single space
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllLiteralString(text, " ")

	return text
}

func printSearchResults(results resultSet, term string) {
	resultQuantifier := "result"
	numResults := len(results)
	if numResults != 1 {
		resultQuantifier += "s"
	}
	fmt.Fprintf(out, "%d %s for '%s'\n", len(results), resultQuantifier, term)

	for num, _ := range results {
		comicNum := strconv.Itoa(num)
		comic := loadComicFromFile(comicNum + ".json")
		url := comicUrl(comicNum)
		padding := fmt.Sprintf("%*s", len(url), "=")
		fmt.Fprintf(out, "\n%s\n%s\n%s\n", url, strings.ReplaceAll(padding, " ", "="), comic.Transcript)
	}
}

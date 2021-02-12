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
	"strings"
)

const baseURL = "https://xkcd.com"
const jsonPath = "info.0.json"
const dataDir = "data/"

type XKCDComic struct {
	Num        int
	Title      string
	Transcript string
}

func main() {
	downloadFlag := flag.Bool("d", false, "Download comics data")

	flag.Parse()

	if *downloadFlag {
		comicNum := "568"
		url := strings.Join([]string{baseURL, comicNum, jsonPath}, "/")
		saveLoc := dataDir + comicNum + ".json"
		downloadComic(url, saveLoc)
	}
	index := buildSearchIndex()

	term := os.Args[len(os.Args)-1]
	search(term, index)
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
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Failed to download comic at url#%s: %s", url, resp.Status)
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
		data, err := ioutil.ReadFile(dataDir + file.Name())
		if err != nil {
			log.Fatalf("Failed to open %v", err)
		}

		var comic XKCDComic
		if err := json.Unmarshal(data, &comic); err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}

		searchableText := comic.Title + comic.Transcript
		for _, word := range strings.Split(searchableText, " ") {
			if index[word] == nil {
				index[word] = make(map[int]bool)
			}
			index[word][comic.Num] = true
		}
	}

	return index
}

func search(term string, index map[string]map[int]bool) {
	if index[term] == nil {
		fmt.Printf("Search term: '%s' not found.\n", term)
		return
	}

	for comicNum, _ := range index[term] {
		fmt.Printf("%v\n", comicNum)
	}
}

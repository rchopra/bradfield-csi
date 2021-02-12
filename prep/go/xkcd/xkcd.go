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

const xkcdURL = "https://xkcd.com/"
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
		downloadComic("570")
	}
	index := buildSearchIndex()

	term := os.Args[len(os.Args)-1]
	search(term, index)
}

func downloadComic(comicNum string) error {
	resp, err := http.Get(xkcdURL + comicNum + "/info.0.json")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("Failed to download comic #%s: %s", comicNum, resp.Status)
	}

	out, err := os.Create(dataDir + comicNum + ".json")
	io.Copy(out, resp.Body)

	return nil
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

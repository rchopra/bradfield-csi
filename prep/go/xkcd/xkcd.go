package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const xkcdURL = "https://xkcd.com/"

type XKCDComic struct {
	Num        int
	Title      string
	Transcript string
}

func main() {
	downloadComic("570")
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

	out, err := os.Create("data/" + comicNum + ".json")
	io.Copy(out, resp.Body)

	return nil
}

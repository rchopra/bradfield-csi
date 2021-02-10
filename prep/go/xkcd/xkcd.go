package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type XKCDComic struct {
	Num        int
	Title      string
	Transcript string
}

func main() {
	contents, err := ioutil.ReadFile("comics.json")
	if err != nil {
		log.Fatalf("ioutil.ReadFile failed: %v", err)
	}

	var comic XKCDComic
	if err := json.Unmarshal(contents, &comic); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	fmt.Printf("Comic: \n%v", comic)
}

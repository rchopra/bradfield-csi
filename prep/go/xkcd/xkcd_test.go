package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	// Supress standard out by having it write to this buffer instead
	out = new(bytes.Buffer)

	var emptyIndex = make(searchIndex)
	var smallIndex = make(searchIndex)
	smallIndex["this"] = resultSet{1: true, 2: true}
	smallIndex["is"] = resultSet{1: true, 2: true}
	smallIndex["a"] = resultSet{1: true}
	smallIndex["test"] = resultSet{1: true, 2: true}

	var tests = []struct {
		term  string
		index searchIndex
		want  int
	}{
		{"anything", emptyIndex, 0},
		{"this", smallIndex, 2},
		{"this!", smallIndex, 2},
		{"THIS", smallIndex, 2},
		{"th is", smallIndex, 0},
		{"is", smallIndex, 2},
		{"a", smallIndex, 1},
		{"test", smallIndex, 2},
		{"bad", smallIndex, 0},
	}

	for _, test := range tests {
		if got := search(test.term, test.index); len(got) != test.want {
			t.Errorf(
				"search(%s, %#v) returned %v results, want %v",
				test.term,
				test.index,
				len(got),
				test.want,
			)
		}
	}
}

func TestBuildSearchIndex(t *testing.T) {
	// This test is going to create two fake xkcd JSON files in a temp data
	// directory and then build the search index out of them. We want the index
	// to normalize and clean the text properly, handling non-alpha characters,
	// whitespace, etc.
	dir, _ := ioutil.TempDir("", "data")
	defer os.RemoveAll(dir)

	f1, _ := ioutil.TempFile(dir, "*.json")
	f2, _ := ioutil.TempFile(dir, "*.json")
	defer os.Remove(f1.Name())
	defer os.Remove(f2.Name())

	num1 := rand.Intn(1000)
	num2 := num1 + 1000
	f1.Write(encodeFakeComic(num1, "the first TITLE", "Lorem  i'psum!"))
	f2.Write(encodeFakeComic(num2, "The 2nd title", "Lorem ipsum\ndolor\nsit amet."))

	only1 := make(resultSet)
	only2 := make(resultSet)
	both := make(resultSet)
	only1[num1] = true
	only2[num2] = true
	both[num1] = true
	both[num2] = true

	tests := []struct {
		token   string
		results resultSet
	}{
		{"the", both},
		{"first", only1},
		{"title", both},
		{"2nd", only2},
		{"lorem", both},
		{"ipsum", both},
		{"dolor", only2},
		{"sit", only2},
		{"amet", only2},
	}

	// Setup complete! Finally exercise this function.
	index := buildSearchIndex(dir + "/")

	// First check that the number of tokens is what we expect
	if len(tests) != len(index) {
		t.Errorf(
			"buildSearchIndex(dataDir) should return a map of length '%d', but got length '%d'\n%v",
			len(tests),
			len(index),
			index,
		)
	}

	// Now check that all tokens we expect are in the index, and for each found
	// token, that the resultSet contains the correct comics (either comic1,
	// comic2, or both).
	for _, test := range tests {
		if _, ok := index[test.token]; !ok {
			t.Errorf(
				"buildSearchIndex(dataDir) should return map containing key '%s' but got %v\n",
				test.token,
				index,
			)
		} else {
			// Since the resultSets are just sets of ints (map[int]bool), DeepEqual
			// is valid as it checks whether the keys map to deeply equal values
			if !reflect.DeepEqual(test.results, index[test.token]) {
				t.Errorf(
					"buildSearchIndex(dataDir) should have resultSet '%v' for token '%s', but got '%v'",
					test.results,
					test.token,
					index[test.token],
				)
			}
		}
	}
}

func TestPrintSearchResults(t *testing.T) {
	// We want to capture standard out and read off it for this test
	out = new(bytes.Buffer)

	dir, _ := ioutil.TempDir("", "data")
	defer os.RemoveAll(dir)

	f, _ := ioutil.TempFile(dir, "*.json")
	defer os.Remove(f.Name())

	re := regexp.MustCompile(`/(\d+)\.json`)
	num, _ := strconv.Atoi(re.FindStringSubmatch(f.Name())[1])
	transcript := "This is an xkcd Comic."
	f.Write(encodeFakeComic(num, "A Comic", transcript))

	results := make(resultSet)
	results[num] = true

	printSearchResults(results, "xkcd", dir+"/")

	// With the result printed to standard out, we want to assert the various
	// parts of the output are appearing correctly
	got := out.(*bytes.Buffer).String()
	gotLines := strings.Split(got, "\n")
	descr := "printSearchResults(results, 'xkcd', dir)"

	// Check the first line contains the number of results
	wantResults := "1 result"
	if !strings.Contains(gotLines[0], "1 result") {
		t.Errorf("%s first line does not match.\n\nGOT:\n%s\n\nWANT:\n%s",
			descr,
			gotLines[0],
			wantResults,
		)
	}

	// Check the comic URL was printed
	wantUrl := comicUrl(strconv.Itoa(num))
	if !strings.Contains(got, wantUrl) {
		t.Errorf(
			"%s URL does not match.\n\nGOT:\n%s\n\nWANT URL:\n%s",
			descr,
			got,
			wantUrl,
		)
	}

	// Check that the transcript was printed
	if !strings.Contains(got, transcript) {
		t.Errorf("%s transcript does not match.\n\nGOT:\n%s\n\nWANT TRANSCRIPT:\n%s",
			descr,
			got,
			transcript,
		)
	}
}

func TestDownloadAllComics(t *testing.T) {
}

// Helper Functions
func encodeFakeComic(num int, title string, transcript string) (data []byte) {
	comic := Comic{Num: num, Title: title, Transcript: transcript}

	data, _ = json.Marshal(comic)
	return data
}

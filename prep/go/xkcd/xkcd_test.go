package main

import "testing"

func TestSearch(t *testing.T) {
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
			t.Errorf("search(%s, %#v) returned %v results, want %v", test.term, test.index, len(got), test.want)
		}
	}
}

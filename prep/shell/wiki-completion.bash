#!/bin/bash

WIKI_BASE="https://en.wikipedia.org/w/api.php"
PARAMS="action=opensearch&format=json&formatversion=2&namespace=0&limit=1"

search_wikipedia() {
  curl -s -XGET "$WIKI_BASE?$PARAMS&search=${COMP_WORDS[1]}" |
    cut -d, -f2 |
    sed -n 's/\[\(.*\)\]/\1/p' |
    sed 's/ /_/g' |
    sed 's/,/ /g'
}

page_completions() {
  OPTIONS=$(search_wikipedia)
  COMPREPLY=($(compgen -W "$OPTIONS" -- "${COMP_WORDS[1]}"))
}

complete -F page_completions wiki

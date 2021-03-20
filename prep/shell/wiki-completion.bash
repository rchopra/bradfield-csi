#!/bin/bash

WIKI_BASE="https://en.wikipedia.org/w/api.php"
PARAMS="action=opensearch&format=json&formatversion=2&namespace=0&limit=1"

search_wikipedia() {
  curl -s -XGET "$WIKI_BASE?$PARAMS&search=${COMP_WORDS[1]}" |
    cut -d, -f2 |                # Cut out the results array, will give "[...]"
    sed -n 's/\[\(.*\)\]/\1/p' | # Remove square brackets
    sed 's/ /_/g' |              # Replace spaces with underscores
    sed 's/,/ /g'                # Replace commas with spaces
}

page_completions() {
  OPTIONS=$(search_wikipedia)
  COMPREPLY=($(compgen -W "$OPTIONS" -- "${COMP_WORDS[1]}"))

  # For section completion, can grab sections from the parse API
  # Need to enable only when the second argument is present, using the first
  # argument as an input
}

complete -F page_completions wiki

#!/bin/bash

# Grabbing JSON from the [MediaWiki
# API](https://www.mediawiki.org/wiki/API:Main_page)
WIKI_BASE="https://en.wikipedia.org/w/api.php"
PARAMS="action=query&format=json&prop=extracts&explaintext=1&redirects=1"

# This is some very basic awk to pull out the section headers from the
# response. The headers look like == Name of Section ==, so I'm pulling out the
# text between the equal signs
function extract_sections() {
  echo "$1" | awk 'BEGIN {
    OFS="\n";
    pattern="n== [^=]* ==";
  }

  {
    while (match($0, pattern)) {
      printf substr($0, RSTART+4, RLENGTH-6) OFS;
      $0=substr($0, RSTART+RLENGTH)
    }
  }'
}

#TODO: Error handling from curl
#TODO: What happens if a page is not found
#TODO: Handle multi-word/punctuation inputs
RESP=$(curl -s -XGET "$WIKI_BASE?$PARAMS&titles=$1")

# I really want to use [jq](https://stedolan.github.io/jq/) for this but it
# seems counter to the spirit of this exercise, so here's some gnarly sed
# instead.
FULL_TEXT=$(echo "$RESP" | sed -n 's/^.*extract":"\(.*\)}}}}.*$/\1/p')

# This is surely buggy, but just finding the first period in the text of the
# article (and adding the period back on the end.
INTRO=${FULL_TEXT%%.*}.
SECTIONS=$(extract_sections "$FULL_TEXT")

printf "%s\n\n%s\n" "$INTRO" "$SECTIONS"

#!/bin/bash

# Grabbing JSON from the [MediaWiki
# API](https://www.mediawiki.org/wiki/API:Main_page)
WIKI_BASE="https://en.wikipedia.org/w/api.php"
PARAMS="action=query&format=json&prop=extracts&explaintext=1&redirects=1"

# This is some very basic awk to pull out the section headers from the
# response. The headers look like == Name of Section ==, so I'm pulling out the
# text between the equal signs. The second argument is the number of equal signs
function extract_sections() {
  echo "$1" | awk "BEGIN {
    OFS=\"\n\";
    pattern=\"n$2 [^=]* $2\";
  }

  {
    while (match(\$0, pattern)) {
      printf substr(\$0, RSTART+${#2}+2, RLENGTH-$((${#2}*2 + 2))) OFS;
      \$0=substr(\$0, RSTART+RLENGTH)
    }
  }"
}

# This is surely buggy, but just taking up to the first period in the text
# passed in and adding the period back on the end.
function extract_intro() {
  echo "${1%%.*}."
}

# Get the entire text of the requested subsection
function section_text() {
  echo "$1" | grep -Eio "== $2 ==\\\n\\\n.*?\\\n\\\n\\\n== "
}

#TODO: Error handling from curl
#TODO: What happens if a page is not found
#TODO: Handle multi-word/punctuation inputs
RESP=$(curl -s -XGET "$WIKI_BASE?$PARAMS&titles=$1")
#RESP=$(cat test.json)

# I really want to use [jq](https://stedolan.github.io/jq/) for this but it
# seems counter to the spirit of this exercise, so here's some gnarly sed
# instead.
FULL_TEXT=$(echo "$RESP" | sed -n 's/^.*extract":"\(.*\)}}}}.*$/\1/p')

INTRO=""
SECTIONS=""

if [[ "$#" -eq 1 ]]; then
  INTRO=$(extract_intro "$FULL_TEXT")
  SECTIONS=$(extract_sections "$FULL_TEXT" "==")
elif [[ "$#" -eq 2 ]]; then
  SECTION_TEXT=$(section_text "$FULL_TEXT" "$2")
  INTRO=$(extract_intro "$(echo "$SECTION_TEXT" | cut -d= -f5 | cut -c 5-)")
  SECTIONS=$(extract_sections "$SECTION_TEXT" "===")
else
  echo "Wrong number of arguments."
  exit 1
fi

printf "%s\n\n%s\n" "$INTRO" "$SECTIONS"

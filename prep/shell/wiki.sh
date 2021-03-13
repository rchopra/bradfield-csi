#!/bin/bash

# Grabbing JSON from the [MediaWiki TextExtracts
# API](https://www.mediawiki.org/wiki/Extension:TextExtracts#API)
WIKI_BASE="https://en.wikipedia.org/w/api.php"
PARAMS="action=query&format=json&prop=extracts&explaintext=1&exsentences=1"

#TODO: Error handling from curl
#TODO: What happens if a page is not found
#TODO: Handle multi-word/punctuation inputs
RESP=$(curl -s -XGET "$WIKI_BASE?$PARAMS&titles=$1")

# I really want to use [jq](https://stedolan.github.io/jq/) for this but it
# seems counter to the spirit of this exercise, so here's some gnarly sed
# instead.
INTRO=$(echo "$RESP" | sed -n 's/^.*extract":"\(.*\)".*$/\1/p')

echo "$INTRO"

# XKCD Program
As described in *[The Go Programming Language](https://www.gopl.io/)* Exercise
4.12.

## Building
In the project directory
```sh
$ go build
```

## Usage
The most basic usage is
```sh
$ ./xkcd term
```
This will look for all the comics downloaded locally in a directory called
`data/`, build a search index, and check if the term exists.

You can also specify optional flags:
- `-d` will attempt to download all missing comics to your save location, which
defaultis to `data/`. Downloads start from the most recent and proceed in
reverse chronological order.
- `-l /save/data/path` will change the download and search paths to the one
specified here.

### Implementation Details
#### Downloading
The program only downloads files when specifically asked (to avoid accidentally
hitting the server), and will only download a file if it is not saved on the
local file system (i.e. in `data/` or as specified with the `-l` flag).

When downloading, the program first tries to determine what the most recent
comic published was and works backwards from there, checking the local file
system for any existing comics and skipping the download if so. This means that
you won't get an updated comic unless you delete the local copy (or specicify a
new save location).

#### Searching
The search index is extremely basic. It's simply a map of strings to sets of
integers (representing the set of comic numbers where that string appears).
There is some very minimal processing done to remove case-sensitivity, strip
punctuation and symbols, and get rid of meta lines like "title text". The
transcripts are tokenized by splitting on spaces.

Multi-word terms are not supported. Any arguments after the first term are
simply dropped.

#### Testing
I'm almost certain this is not the right way to test, and there's so much more
that can and should be done (I'm not testing all the error cases, for one), but
I found it somewhat of a challenge to write these -- I'm so used to all the
conveniences of languages with test huge test frameworks.

### Extensions
I'd love to make the downloading concurrent, but I feel like I've already spent
a bit too much time on this and need to pick up the pace on other parts of the
prep. Will definitely revisit though.

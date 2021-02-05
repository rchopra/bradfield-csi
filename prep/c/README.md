## A minimal `ls` clone

My goal was to stay as faithful to the version of `ls` I'm running on my
development machine (MacOS) as possible. This will work with the current
directory by default or any path you give it. It does not work if you give it a
file, which I've left as a TODO for now in the interest of continuing to move
forward with the rest of the prep work.

### Flags implemented
* `-a`: Shows files that start with `'.'`
* `-s`: Shows the total size in the directory along with the size of each entry
        in blocks.
* `-k`: Same as `-s`, but in kilobytes. On my system this is half the number of
        blocks -- the block size is 512 bytes.
* `-S`: Sort by number file size in bytes, descending.
* `-1`: Prints all listings on a new line

### Open questions
  * I wasn't sure how to support sorting without collecting everything into an
    array first, so I had to pick some arbitrary number (2^16 - 1) for files in
    a directory, but I'm sure there's something more intelligent to do.
  * Style wise, I'm not sure if I decomposed too much or not enough. I'm used to
    writing five-line methods in Ruby.
  * I thought a bit about how to elegantly handle when the argument is a file
    instead of directory without duplicating a bunch of code, but it was not
    immediately obvious.

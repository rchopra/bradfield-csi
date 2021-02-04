#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/stat.h>

#define DEFAULT_DIR "."
#define MAX_ENTRIES 65535
#define SHOW_HIDDEN 0b01
#define SHOW_SIZE   0b10

typedef struct lsent {
  char name[1024];
  blkcnt_t blocks;
  off_t size;
} LSENT;

void print_entry(LSENT *entry, unsigned char flags) {
  if (flags & SHOW_SIZE) {
    printf("%3llu %-9s\n", entry->blocks, entry->name);
  } else {
    printf("%-9s\n", entry->name);
  }
}

int collect_contents(DIR *dirp, LSENT *entries, unsigned char flags) {
  // TODO: Error checking when calling `readdir`
  struct dirent *direntp;
  int dir_fd = dirfd(dirp);
  LSENT *start = entries;

  while((direntp = readdir(dirp))) {
    struct stat filestats;

    if (fstatat(dir_fd, direntp->d_name, &filestats, 0)) {
      printf("%s\n", strerror(errno));
      exit(1);
    }

    if (direntp->d_name[0] == '.' && !(flags & SHOW_HIDDEN)) {
      continue;
    }

    strcpy(entries->name, direntp->d_name);
    entries->blocks = filestats.st_blocks;
    entries->size   = filestats.st_size;
    entries++;
  }

  return entries - start;
}

int main(int argc, char *argv[]) {
  unsigned char flags = 0;
  char *path;
  char c;

  // Flag parsing lifted from K&R C, Section 5.10 (pg. 117)
  while (--argc > 0 && (*++argv)[0] == '-') {
    while ((c = *++argv[0])) {
      switch(c) {
        case 'a':
          flags |= SHOW_HIDDEN;
          break;
        case 's':
          flags |= SHOW_SIZE;
          break;
        default:
          printf("myls: illegal option %c\n", c);
          argc = 0;
          break;
      }
    }
  }

  path = argc == 0 ? DEFAULT_DIR : *argv;

  DIR *dirp = opendir(path);
  LSENT *entries = malloc(MAX_ENTRIES * sizeof(LSENT));
  int num_entries;

  // TODO: Handle when input is just a file
  if (dirp == NULL) {
    printf("myls: %s: %s\n", path, strerror(errno));
    return 1;
  }

  // TODO: Sort entries like `ls` by default
  num_entries = collect_contents(dirp, entries, flags);
  for (int i = 0; i < num_entries; i++) {
    print_entry(&entries[i], flags);
  }

  free(entries);
  closedir(dirp);
  return 0;
}

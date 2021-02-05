#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/stat.h>

#define MAX_DIR_ENTRIES 65535

// Bit masks for ls flags
#define SHOW_HIDDEN    0b00001
#define SHOW_SIZE      0b00010
#define SHOW_KILOBYTES 0b00100
#define NEW_LINES      0b01000
#define SORT_BY_SIZE   0b10000

typedef struct lsent {
  char name[1024];
  blkcnt_t blocks;
  off_t bytes;
} LSENT;

struct ls_stats {
  blkcnt_t total_blocks;
  off_t total_size;
};

// Function prototypes
void print_summary(struct ls_stats *summary, unsigned char flags);
void print_entry(LSENT *entry, unsigned char flags);
int compare_by_size(const void *first, const void *second);
int collect_contents(DIR *dirp, LSENT *entries, struct ls_stats *summary, unsigned char flags);

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
        case 'k':
          flags |= SHOW_KILOBYTES;
          break;
        case 's':
          flags |= SHOW_SIZE;
          break;
        case 'S':
          flags |= SORT_BY_SIZE;
          break;
        case '1':
          flags |= NEW_LINES;
          break;
        default:
          printf("myls: illegal option %c\n", c);
          argc = 0;
          break;
      }
    }
  }

  path = argc == 0 ? "." : *argv;

  DIR *dirp = opendir(path);
  LSENT *entries = malloc(MAX_DIR_ENTRIES * sizeof(LSENT));
  int num_entries;

  // TODO: Handle when input is just a file
  if (dirp == NULL) {
    printf("myls: %s: %s\n", path, strerror(errno));
    return 1;
  }

  struct ls_stats summary = { 0, 0 };

  num_entries = collect_contents(dirp, entries, &summary, flags);
  if (flags & SORT_BY_SIZE) {
    qsort(entries, num_entries, sizeof(LSENT), compare_by_size);
  }

  print_summary(&summary, flags);
  for (int i = 0; i < num_entries; i++) {
    print_entry(&entries[i], flags);
  }
  if (!(flags & NEW_LINES)) {
    printf("\n");
  }

  free(entries);
  closedir(dirp);
  return 0;
}

int collect_contents(DIR *dirp, LSENT *entries, struct ls_stats *summary, unsigned char flags) {
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
    entries->blocks        = filestats.st_blocks;
    entries->bytes         = filestats.st_size;
    summary->total_blocks += entries->blocks;
    summary->total_size   += entries->bytes;
    entries++;
  }

  return entries - start;
}

void print_summary(struct ls_stats *summary, unsigned char flags) {
  if (flags & SHOW_SIZE) {
    int size_scale = (flags & SHOW_KILOBYTES) ? 2 : 1;
    printf("total %llu\n", summary->total_blocks / size_scale);
  }
}

void print_entry(LSENT *entry, unsigned char flags) {
  char sep = (flags & NEW_LINES) ? '\n' : ' ';
  if (flags & SHOW_SIZE) {
    int size_scale = (flags & SHOW_KILOBYTES) ? 2 : 1;
    printf("%3llu %-9s%c", entry->blocks / size_scale, entry->name, sep);
  } else {
    printf("%-9s%c", entry->name, sep);
  }
}

int compare_by_size(const void *first, const void *second) {
  off_t fsize = ((LSENT *)first)->bytes;
  off_t ssize = ((LSENT *)second)->bytes;
  return (fsize < ssize) - (fsize > ssize);
}

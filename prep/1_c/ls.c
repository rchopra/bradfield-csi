#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/stat.h>

#define DEFAULT_DIR "."
#define MAX_ENTRIES 65535

typedef struct lsent {
  char name[1024];
  blkcnt_t blocks;
  off_t size;
} LSENT;

void print_entry(LSENT *entry) {
  printf("%3llu %-9s\n", entry->blocks, entry->name);
}

int collect_contents(DIR *dirp, LSENT *entries) {
  // TODO: Implement the -a flag; it is on by default now
  // TODO: Error checking when calling `readdir`
  struct dirent *direntp;

  LSENT *start = entries;
  while((direntp = readdir(dirp))) {
    struct stat filestats;

    if (stat(direntp->d_name, &filestats)) {
      printf("%s\n", strerror(errno));
      exit(1);
    }

    strcpy(entries->name, direntp->d_name);
    entries->blocks = filestats.st_blocks;
    entries->size   = filestats.st_size;
    entries++;
  }

  return entries - start;
}

int main(int argc, char *argv[]) {
  // TODO: This will have to change when introducing optional flags
  char *path = (argc == 1) ? DEFAULT_DIR : argv[1];
  DIR *dirp = opendir(path);
  LSENT *entries = malloc(MAX_ENTRIES * sizeof(LSENT));
  int num_entries;

  // TODO: Handle when input is just a file
  if (dirp == NULL) {
    printf("%s: %s: %s\n", argv[0], path, strerror(errno));
    return 1;
  }

  // TODO: Sort entries like `ls` by default
  num_entries = collect_contents(dirp, entries);
  for (int i = 0; i < num_entries; i++) {
    print_entry(&entries[i]);
  }

  free(entries);
  closedir(dirp);
  return 0;
}

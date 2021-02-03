#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/stat.h>

#define DEFAULT_DIR "."

void print_entry(struct dirent *ent) {
  struct stat filestats;

  if (stat(ent->d_name, &filestats)) {
    printf("%s\n", strerror(errno));
    exit(1);
  }
  printf("%3llu %-9s\n", filestats.st_blocks, ent->d_name);
}

void list_contents(DIR *dirp) {
  // TODO: Sort files like `ls` by default
  // TODO: Implement the -a flag; it is on by default now
  struct dirent *ent;

  // TODO: Error checking when calling `readdir`
  while((ent = readdir(dirp))) {
    print_entry(ent);
  }
}

int main(int argc, char *argv[]) {
  // TODO: This will have to change when introducing optional flags
  char *path = (argc == 1) ? DEFAULT_DIR : argv[1];
  DIR *dirp = opendir(path);

  // TODO: Handle when input is just a file
  if (dirp == NULL) {
    printf("%s: %s: %s\n", argv[0], path, strerror(errno));
    return 1;
  }

  list_contents(dirp);

  closedir(dirp);
  return 0;
}

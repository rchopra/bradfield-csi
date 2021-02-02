#include <dirent.h>
#include <fcntl.h>
#include <stdio.h>

#define DEFAULT_DIR "."

void list_contents(DIR *dirp) {
  // TODO: Sort files like `ls` by default
  // TODO: Implement the -a flag; it is on by default now
  struct dirent *ent;
  // TODO: Error checking when calling `readdir`
  while((ent = readdir(dirp))) {
    //TODO: `ls` uses some kind of padding, not tab characters
    printf("%s\t", ent->d_name);
  }
  printf("\n");
}

int main(int argc, char *argv[]) {
  // TODO: This will have to change when introducing optional flags
  char *path = (argc == 1) ? DEFAULT_DIR : argv[1];
  DIR *dirp = opendir(path);
  // TODO: Some error checking

  list_contents(dirp);

  closedir(dirp);
  return 0;
}

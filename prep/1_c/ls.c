#include <dirent.h>
#include <fcntl.h>
#include <stdio.h>

void list_contents(DIR *dirp) {
  // TODO: Sort files like `ls` by default
  // TODO: Implement the -a flag; it is on by default now
  struct dirent *ent;
  // TODO: Error checking when calling `readdir`
  while((ent = readdir(dirp))) {
    printf("%s\t", ent->d_name);
  }
  printf("\n");
}

int main(int argc, char *argv[]) {
  // TODO: Read the directory from the command line
  char *path = ".";
  DIR *dirp;

  dirp = opendir(path);
  // TODO: Some error checking

  list_contents(dirp);

  closedir(dirp);
  return 0;
}

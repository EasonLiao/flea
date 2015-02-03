package builtin

import (
  "fmt"
  "github.com/easonliao/flea/core"
  "os"
  "path/filepath"
)

// Converts tree path to FS path relative to repository path.
// For example, if the treePath is /src/helloworld.go, it will be converted
// to src/helloworld.go in Linux.
func TreePathToRelFsPath(treePath string) string {
  return filepath.FromSlash(treePath[1:])
}

// Converts path relative to current working directory to path relative to root path of repository.
// For example, if the base path of repo is in /home/user/rep, and you are in /home/user/rep/src,
// the relPath is golang/hello.go, this function will return src/golang/hello.go
func GetRelFsPath(relPath string) string {
  cwd, err := os.Getwd()
  if err != nil {
    panic("Error in getting working directory.")
  }
  if path, err := filepath.Rel(core.GetRepoDirectory(), cwd); err != nil {
    panic(err.Error())
  } else {
    return filepath.Join(path, relPath)
  }
}

// Converts the path relative to repo to standard tree path.
func RelFsPathToTreePath(fsPath string) string {
  return "/" + filepath.ToSlash(fsPath)
}

// Print the str and exit.
func PrintAndExit(str string) {
  fmt.Println(str)
  os.Exit(1)
}

package builtin

import (
  "log"
  "os"
  "path"
)

func CmdInit() int {
  // Get current working directory.
  cwd, err := os.Getwd()
  if err != nil {
    log.Fatal("Failed to get current working directory.")
  }
  fleaDir := path.Join(cwd, ".flea")
  if _, err := os.Stat(fleaDir); err == nil {
    log.Fatal(".flea had already existed in current directory.")
  } else if ! os.IsNotExist(err) {
    log.Fatal("Error", err)
  }
  initDirectory(fleaDir)
  return 0
}

func initDirectory(baseDir string) {
  os.Mkdir(baseDir, os.ModeDir | 0777)
  os.Mkdir(path.Join(baseDir, "objects"), os.ModeDir | 0777)
  os.Mkdir(path.Join(baseDir, "refs"), os.ModeDir | 0777)
  os.Mkdir(path.Join(baseDir, "infos"), os.ModeDir | 0777)
}

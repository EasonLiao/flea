package setup

import (
  "github.com/easonliao/flea/core"
  "log"
  "os"
  "path"
)

var FleaDir string
var isSetup bool = false

func SetupFleaDir() {
  cwd, err := os.Getwd()
  if err != nil {
    log.Fatal("Failed in getting current directory.")
  }
  // Find flea directory.
  curDir := cwd
  for {
    if _, err := os.Stat(path.Join(curDir, ".flea")); err != nil {
      prevDir := curDir
      curDir = path.Dir(curDir)
      if prevDir == curDir {
        log.Fatal("Can't find .flea directory in current and parent directories.")
      }
    } else {
      FleaDir = path.Join(curDir, ".flea")
      isSetup = true
      break
    }
  }
  core.InitStoreDir(path.Join(FleaDir, "objects"))
}

func AssertIsSetup() {
  if !isSetup {
    panic("Flea directory has not been setup.")
  }
}

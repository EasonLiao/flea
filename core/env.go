package core

import (
  "errors"
  "log"
  "os"
  "path"
)

var (
  ErrNoFleaDir = errors.New("core: can't find flea directory")
  ErrFleaDirExist = errors.New("core: .flea dir has already existed in current directory.")
)

var (
  initialized = false
  workingDirectory = ""
  fleaDirectory = ""
  storeDirectory = ""
)

func initPaths(wd string) {
  workingDirectory = wd
  fleaDirectory = path.Join(workingDirectory, ".flea")
  storeDirectory = path.Join(fleaDirectory, "objects")
  log.Println(workingDirectory, fleaDirectory, storeDirectory)
  initialized = true
}

// Creating a new Flea repository in current working directory.
func InitNew() error {
  // Get current working directory.
  cwd, err := os.Getwd()
  if err != nil {
    return err
  }
  fd := path.Join(cwd, ".flea")
  if _, err := os.Stat(fd); err == nil {
    err = ErrFleaDirExist
    return err
  } else if ! os.IsNotExist(err) {
    return err
  }
  os.Mkdir(fd, os.ModeDir | 0777)
  os.Mkdir(path.Join(fd, "objects"), os.ModeDir | 0777)
  os.Mkdir(path.Join(fd, "refs"), os.ModeDir | 0777)
  os.Mkdir(path.Join(fd, "infos"), os.ModeDir | 0777)
  initPaths(cwd)
  return nil
}

// Initializing Flea from an existing Flea repository.
func InitFromExisting() error {
  cwd, err := os.Getwd()
  if err != nil {
    return err
  }
  curDir := cwd
  for {
    if _, err := os.Stat(path.Join(curDir, ".flea")); err != nil {
      prevDir := curDir
      curDir = path.Dir(curDir)
      if prevDir == curDir {
        return ErrNoFleaDir
      }
    } else {
      initPaths(cwd)
      break
    }
  }
  return nil
}

// Get the root working directory of current Flea repository.
func GetWorkingDirectory() string {
  assertInit()
  return workingDirectory
}

// Get the full path of .flea directory of current Flea repository.
func GetFleaDirectory() string {
  assertInit()
  return fleaDirectory
}

// Get the full path of .flea/objects directory of current Flea repository.
func GetStoreDirectory() string {
  assertInit()
  return storeDirectory
}

func assertInit() {
  if !initialized {
    panic("Core package has not been initialized.")
  }
}

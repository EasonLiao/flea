package core

import (
  "errors"
  "log"
  "os"
  "path/filepath"
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
  pathPrefix = ""
)

func initPaths(wd string) {
  cd, _ := os.Getwd()
  workingDirectory = wd
  fleaDirectory = filepath.Join(workingDirectory, ".flea")
  storeDirectory = filepath.Join(fleaDirectory, "objects")
  pathPrefix, _ = filepath.Rel(workingDirectory, cd)
  if pathPrefix == "." {
    pathPrefix = ""
  }
  pathPrefix = "/" + pathPrefix
  log.Println(workingDirectory, fleaDirectory, storeDirectory, pathPrefix)
  initialized = true
}

// Creating a new Flea repository in current working directory.
func InitNew() error {
  // Get current working directory.
  cwd, err := os.Getwd()
  if err != nil {
    return err
  }
  fd := filepath.Join(cwd, ".flea")
  if _, err := os.Stat(fd); err == nil {
    err = ErrFleaDirExist
    return err
  } else if ! os.IsNotExist(err) {
    return err
  }
  os.Mkdir(fd, os.ModeDir | 0777)
  os.Mkdir(filepath.Join(fd, "objects"), os.ModeDir | 0777)
  os.Mkdir(filepath.Join(fd, "refs"), os.ModeDir | 0777)
  os.Mkdir(filepath.Join(fd, "infos"), os.ModeDir | 0777)
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
    if _, err := os.Stat(filepath.Join(curDir, ".flea")); err != nil {
      prevDir := curDir
      curDir = filepath.Dir(curDir)
      if prevDir == curDir {
        return ErrNoFleaDir
      }
    } else {
      initPaths(curDir)
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

// Gets the prefix of the path.
func GetPathPrefix() string {
  assertInit()
  return pathPrefix
}

func assertInit() {
  if !initialized {
    panic("Core package has not been initialized.")
  }
}

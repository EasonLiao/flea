package core

import (
  "encoding/hex"
  "errors"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "strings"
)

var (
  ErrNoFleaDir = errors.New("core: can't find flea directory")
  ErrFleaDirExist = errors.New("core: .flea dir has already existed in current directory.")
  ErrNoHeadFile = errors.New("core: HEAD file doesn't exist")
  ErrNotBranch = errors.New("core: not in a valid branch")
)

var (
  initialized = false
  workingDirectory = ""
  fleaDirectory = ""
  storeDirectory = ""
  pathPrefix = ""
  headFilePath = ""
  branchHeadDir = ""
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
  headFilePath = filepath.Join(fleaDirectory, "HEAD")
  branchHeadDir = filepath.Join(fleaDirectory, filepath.Join("refs", "heads"))
  // log.Println(workingDirectory, fleaDirectory, storeDirectory, pathPrefix)
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
  os.Mkdir(filepath.Join(fd, filepath.Join("refs", "heads")), os.ModeDir | 0777)
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

// Get the full path of branch head dir. /refs/heads/
func GetBranchHeadDir() string {
  assertInit()
  return branchHeadDir
}

// Gets current branch. The return values can be 1 of 3:
// 1) branch name and nil.
// 2) empty branch and ErrNoHeadFile.
// 3) empty branch and ErrNotBranch.
func GetCurrentBranch() (branch string, err error) {
  data, err := ioutil.ReadFile(getHeadFilePath())
  content := string(data)
  if err == nil {
    if strings.HasPrefix(content, "ref:") {
      // It contains a link to branch name.
      branch = content[len("ref:"):]
    } else {
      // It contains a hash value of the commit object.
      err = ErrNotBranch
    }
  } else if os.IsNotExist(err) {
    // There's no HEAD file.
    err = ErrNoHeadFile
  } else {
    log.Fatal("Error in reading file %s", getHeadFilePath())
  }
  return
}

// Gets current position in commit history. The return values can be:
// 1) A valid CommitTree object and nil.
// 2) nil and ErrNoHeadFile.
func GetCurrentCommit() (*Commit, error) {
  branch, err := GetCurrentBranch()
  var commitHash []byte
  if err == nil {
    // Now we're in a valid branch, gets the commit hash from branch file.
    commitHash, err = ioutil.ReadFile(filepath.Join(GetBranchHeadDir(), branch))
    if err != nil {
      log.Fatalf("Can't read the branch file %s.", branch)
    }
  } else if err ==  ErrNotBranch {
    // We're not in a branch, getting the current position from HEAD file.
    commitHash, err = ioutil.ReadFile(getHeadFilePath())
    if err != nil {
      log.Fatal("Can't read the HEAD file.")
    }
  } else {
    // The err should be ErrNoHeadFile
    return nil, err
  }
  if commitHash, err = hex.DecodeString(string(commitHash)); err != nil {
    log.Fatal("Not a valid hash string.")
  }
  fType, data, err := GetCAStore().Get(commitHash)
  if err != nil {
    log.Fatal("Failed to get %x from CAStore", commitHash)
  }
  if fType != CommitType {
    log.Fatal("The hash %x doesn't point to a commit object.", commitHash)
  }
  commit, err := fromBytesToCommitObject(data)
  if err != nil {
    log.Fatal("Failed to convert file in %x to commit object.", commitHash)
  }
  return commit, nil
}

func WriteHeadFile(data []byte) {
  if err := ioutil.WriteFile(getHeadFilePath(), data, 0777); err != nil {
    panic("Failed to write to HEAD file.")
  }
}

// Checks whether a branch is valid or not.
func IsValidBranch(branch string) bool {
  if _, err := os.Stat(filepath.Join(GetBranchHeadDir(), branch)); err == nil {
    return true
  }
  return false
}

func getHeadFilePath() string {
  assertInit()
  return headFilePath
}

func assertInit() {
  if !initialized {
    panic("Core package has not been initialized.")
  }
}

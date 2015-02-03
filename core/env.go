package core

import (
  "encoding/hex"
  "errors"
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
  ErrInvalidBranch = errors.New("core: Invalid branch")
)

var (
  initialized = false
  repoDirectory = ""
  fleaDirectory = ""
  storeDirectory = ""
  pathPrefix = ""
  headFilePath = ""
  branchHeadDir = ""
)

func initPaths(wd string) {
  cd, _ := os.Getwd()
  repoDirectory = wd
  fleaDirectory = filepath.Join(repoDirectory, ".flea")
  storeDirectory = filepath.Join(fleaDirectory, "objects")
  pathPrefix, _ = filepath.Rel(repoDirectory, cd)
  if pathPrefix == "." {
    pathPrefix = ""
  }
  pathPrefix = "/" + pathPrefix
  headFilePath = filepath.Join(fleaDirectory, "HEAD")
  branchHeadDir = filepath.Join(fleaDirectory, filepath.Join("refs", "heads"))
  // log.Println(repoDirectory, fleaDirectory, storeDirectory, pathPrefix)
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
  if exists(fd) {
    return ErrFleaDirExist
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
    panic(err.Error())
  }
  curDir := cwd
  for {
    if !exists(filepath.Join(curDir, ".flea")) {
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

// Get the root directory of current Flea repository.
func GetRepoDirectory() string {
  assertInit()
  return repoDirectory
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
  if !exists(getHeadFilePath()) {
    err = ErrNoHeadFile
    return
  }
  data, _ := read(getHeadFilePath())
  content := string(data)
  if strings.HasPrefix(content, "ref:") {
      // It contains a link to branch name.
      branch = content[len("ref:"):]
  } else {
    // It contains a hash value of the commit object.
    err = ErrNotBranch
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
    commitHash, err = read(filepath.Join(GetBranchHeadDir(), branch))
    if err != nil {
      log.Fatalf("Can't read the branch file %s.", branch)
    }
  } else if err ==  ErrNotBranch {
    // We're not in a branch, getting the current position from HEAD file.
    commitHash, err = read(getHeadFilePath())
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

// Gets the hash of the HEAD of a branch.
func GetBranchHead(branch string) []byte {
  if !IsValidBranch(branch) {
    log.Fatalf("%s is not a valid branch.\n", branch)
  }
  head, _ := read(filepath.Join(GetBranchHeadDir(), branch))
  hash, err := hex.DecodeString(string(head))
  if err != nil {
    panic(err.Error())
  }
  if !GetCAStore().Exists(hash) {
    log.Fatalf("%x doesn't exist in repo\n", hash)
  }
  return hash
}

// Updates the head commit of a the branch.
func UpdateBranchHead(branch string, commitHash []byte) {
  if !GetCAStore().Exists(commitHash) {
    panic("Not a valid commit hash.")
  }
  hashString := hex.EncodeToString(commitHash)
  write(filepath.Join(GetBranchHeadDir(), branch), []byte(hashString))
}

// Updates the HEAD file.
func WriteHeadFile(data []byte) {
  write(getHeadFilePath(), data)
}

// Checks whether a branch is valid or not.
func IsValidBranch(branch string) bool {
  if head, err := read(filepath.Join(GetBranchHeadDir(), branch)); err == nil {
    if hash, err := hex.DecodeString(string(head)); err == nil {
      if GetCAStore().Exists(hash) {
        return true
      }
    }
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

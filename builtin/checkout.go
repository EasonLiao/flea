package builtin

import (
  "encoding/hex"
  "fmt"
  "io/ioutil"
  "github.com/easonliao/flea/core"
  "os"
  "path/filepath"
  "sort"
)

func UsageCheckout() string {
  return "flea checkout (<branch>|<commit-hash>)"
}

func CmdCheckout() error {
  if len(os.Args) != 3 {
    PrintAndExit(UsageCheckout())
  }
  dest := os.Args[2]
  var commit *core.Commit

  if core.IsValidBranch(dest) {
    // Checkout to a branch.
    var err error
    head := core.GetBranchHead(dest)
    commit, err = core.GetCommitObject(head)
    if err == core.ErrNoMatch {
      fmt.Println("The head of the branch doesn't point to a valid commit object")
      os.Exit(1)
    }
    fmt.Printf("checking out to %s branch with hash %x\n", dest, head)
    deleteAllFilesInCurrentCommit()
    restoreRepoFromCommit(commit)
    core.WriteHeadFile([]byte("ref:" + dest))
  } else {
    // Checkout to a commit.
    hashPrefix, err := hex.DecodeString(dest)
    if err != nil {
      fmt.Println("Not a valid hash string.")
      os.Exit(1)
    }
    hashs := core.GetCAStore().GetMatchedHashs(hashPrefix)
    if len(hashs) > 1 {
      fmt.Println("More than one object match the hash value.")
      os.Exit(1)
    }
    if len(hashs) == 0 {
      fmt.Println("No matched commit object found.")
      os.Exit(1)
    }
    commit, _ = core.GetCommitObject(hashs[0])
    fmt.Printf("checking out to commit %x\n", commit.GetCommitHash())
    deleteAllFilesInCurrentCommit()
    restoreRepoFromCommit(commit)
    core.WriteHeadFile([]byte(hex.EncodeToString(hashs[0])))
  }
  return nil
}

func deleteAllFilesInCurrentCommit() {
  commit, err := core.GetCurrentCommit()
  if err == core.ErrNoHeadFile {
    // Commit history is empty, nothing to do.
    return
  } else if err != nil {
    panic(err.Error())
  }
  tree := commit.GetCATree()
  paths := make([]string, 0, 64)
  fn := func(treePath string, node core.Node) error {
    paths = append(paths, treePath)
    return nil
  }
  tree.Traverse(fn, "/")
  // Sorts the path in descending order so we'll delete files/dirs in reverse order of
  // the namespace hierarchy.
  sort.Sort(sort.Reverse(sort.StringSlice(paths)))

  for _, path := range(paths) {
    fullPath := filepath.Join(core.GetRepoDirectory(), TreePathToRelFsPath(path))
    // Deletes file.
    os.Remove(fullPath)
  }
}

func restoreRepoFromCommit(commit *core.Commit) {
  restore := func(treePath string, node core.Node) error {
    fsPath := filepath.Join(core.GetRepoDirectory(), TreePathToRelFsPath(treePath))
    if node.IsDir() {
      os.Mkdir(fsPath, 0777)
    } else {
      data, _  := node.GetData()
      err := ioutil.WriteFile(fsPath, data, 0666)
      return err
    }
    return nil
  }
  err := commit.GetCATree().Traverse(restore, "/")
  if err != nil {
    fmt.Println("warning:", err.Error())
  }
}

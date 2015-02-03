package builtin

import (
  "fmt"
  "github.com/easonliao/flea/core"
)

func CmdStatus() error {
  idxTree := core.GetIndexTree()
  fsTree := core.GetFsTree()
  commit, err := core.GetCurrentCommit()
  var commitTree core.Tree
  if err == core.ErrNoHeadFile {
    // We're not in any commit point, the history is empty.
    // Creates an empty MemTree.
    commitTree = core.NewMemTree()
  } else {
    // Gets the tree of current commit.
    commitTree = commit.GetCATree()
  }

  // First compares the commit tree(CATree) to staging area(IndexTree).
  deleted, newFiles, diffes := core.CompareTrees(commitTree, idxTree)
  if len(deleted) > 0 || len(newFiles) > 0 || len(diffes) > 0 {
    fmt.Println("Changes to be committed:\n")
    for _, file := range(deleted) {
      fmt.Printf("\t%s\n", TreePathToRelFsPath(file))
    }
    for _, file := range(newFiles) {
      fmt.Printf("\t%s\n", TreePathToRelFsPath(file))
    }
    for _, file := range(diffes) {
      fmt.Printf("\t%s\n", TreePathToRelFsPath(file))
    }
    fmt.Println("")
  }

  // Compares the staging area(IndexTree) to working directory(FsTree).
  deleted, untracked, diffes := core.CompareTrees(idxTree, fsTree)
  if len(deleted) > 0 || len(diffes) > 0 {
    fmt.Println("Changes not statged for commit:\n")
    for _, file := range(deleted) {
      fmt.Printf("\t%s\n", TreePathToRelFsPath(file))
    }
    for _, file := range(diffes) {
      fmt.Printf("\t%s\n", TreePathToRelFsPath(file))
    }
    fmt.Println("")
  }

  if len(untracked) > 0 {
    fmt.Println("Untracked files:\n")
    for _, file := range(untracked) {
      fmt.Printf("\t%s\n", TreePathToRelFsPath(file))
    }
    fmt.Println("")
  }
  return nil
}

package builtin

import (
  "bytes"
  "flag"
  "fmt"
  "github.com/easonliao/flea/core"
  "os"
  "os/user"
)

func CmdCommit() error {
  if len(os.Args) <= 2 {
    fmt.Println("Not enough arguments.")
    os.Exit(1)
  }
  flags := flag.NewFlagSet("commit", 0)
  comment := flags.String("m", "No Comment", "comment")
  flags.Parse(os.Args[2:])

  indexTree := core.GetIndexTree()
  commit, err := core.GetCurrentCommit()

  if err == nil {
    if bytes.Compare(commit.Tree, indexTree.GetHash()) == 0 {
      // Compares the hash of the commit tree in to the hash of the index tree, if they
      // match then there's nothing to be committed.
      fmt.Println("There's nothing to commit")
      os.Exit(0)
    }
  }
  branch, err := core.GetCurrentBranch()
  if err == core.ErrNotBranch {
    // We're in non-branch, can't commit anything.
    fmt.Printf("Can't commit in a non-branch.")
    os.Exit(1)
  }

  // Creats a CATree from staging area.
  caTree, err := core.BuildCATreeFromIndexFile()
  if err != nil {
    if err == core.ErrEmptyTree {
      fmt.Println("The staging area is empty.")
      os.Exit(1)
    }
    fmt.Printf("Failed to build CATree from staging area: %s\n", err.Error())
    os.Exit(1)
  }

  // Hash of current commit, or nil if there's no commit in history of current branch.
  var commitHash []byte = nil
  if commit != nil {
    commitHash = commit.GetCommitHash()
  }

  var username string = "unknown"
  if user, err := user.Current(); err == nil {
    username = user.Username
  }
  // Creates a commit object.
  hash, err := core.CreateCommitObject(caTree.GetHash(), commitHash, username, *comment)

  if err != nil {
    fmt.Printf("Failed to create the commit object: %s", err.Error())
    os.Exit(1)
  }

  if _, err := core.GetCurrentBranch(); err == nil {
    // We are in a valid branch, just update the HEAD of the branch.
    core.UpdateBranchHead(branch, hash)
  } else if err == core.ErrNoHeadFile {
    // There's no history and branch. Creates a default master branch and updates its HEAD.
    core.WriteHeadFile([]byte("ref:master"))
    branch = "master"
    core.UpdateBranchHead(branch, hash)
  }
  return nil
}

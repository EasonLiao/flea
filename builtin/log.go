package builtin

import (
  "fmt"
  "github.com/easonliao/flea/core"
)

func CmdLog() error {
  commit, err := core.GetCurrentCommit()
  if err == nil {
    for commit != nil {
      printCommit(commit)
      commit = commit.GetPrevCommit()
    }
  }
  return err
}

func printCommit(commit *core.Commit) {
  fmt.Printf("Commit: %x\n", commit.GetCommitHash())
  fmt.Printf("Author: %s\n", commit.Author)
  fmt.Printf("Comment: %s\n", commit.Comment)
  fmt.Println("")
}

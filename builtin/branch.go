package builtin

import (
  "fmt"
  "github.com/easonliao/flea/core"
  "os"
)

func UsageBranch() {
  fmt.Println("Usage: flea branch")
  os.Exit(1)
}

func CmdBranch() error {
  branch, err := core.GetCurrentBranch()
  if err == core.ErrNotBranch {
    fmt.Println("<not on the head of a branch>")
    os.Exit(1)
  } else {
    fmt.Println(branch)
  }
  return nil
}

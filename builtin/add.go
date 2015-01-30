package builtin

import (
  "fmt"
  "github.com/easonliao/flea/core"
  "os"
  "path/filepath"
)

func CmdAdd() error {
  addPath := os.Args[len(os.Args) - 1]

  treePath := filepath.ToSlash(filepath.Join(core.GetPathPrefix(), addPath))
  fsTree := core.GetFsTree()
  idxTree := core.GetIndexTree()

  node, err := fsTree.Get(treePath)

  if err != nil {
    return nil
  }

  _, _ = fsTree, idxTree
  fmt.Println("add", treePath)
  fmt.Printf("%x\n", node.GetHashValue())
  return nil
}

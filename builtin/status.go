package builtin

import (
  "fmt"
  "github.com/easonliao/flea/core"
)

func CmdStatus() error {
  var idxTree *core.IndexTree = core.GetIndexTree()
  var fsTree *core.FsTree = core.GetFsTree()

  fMisses, iMisses, diffes := core.CompareTrees(idxTree, fsTree)
  showStatus(fMisses, iMisses, diffes)
  return nil
}

func showStatus(fsMisses, idxMisses, diffes []string) {
  if len(diffes) > 0 {
    printFiles("Changes not stated for commit:", diffes)
  }
  if len(idxMisses) > 0 {
    printFiles("Untracked files:", idxMisses)
  }
}

func printFiles(title string, files []string) {
  fmt.Printf("%s\n\n", title)
  for _, file := range(files) {
    fmt.Println("\t", file)
  }
  fmt.Println("")
}

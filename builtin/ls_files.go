package builtin

import (
  "fmt"
  "github.com/easonliao/flea/core"
  "os"
  "path/filepath"
)

func UsageLsFiles() {
  fmt.Println("Usage: flea ls-files")
  os.Exit(1)
}

func CmdLsFiles() error {
  if commit, err := core.GetCurrentCommit(); err == nil {
    tree := commit.GetCATree()
    filepaths := make([]string, 0)
    fn := func(treePath string, node core.Node) error {
      if !node.IsDir() {
        filepaths = append(filepaths, treePath)
      }
      return nil
    }
    tree.Traverse(fn, "/")
    if cwd, err := os.Getwd(); err == nil {
      for _, path := range(filepaths) {
        fullPath := filepath.Join(core.GetRepoDirectory(), TreePathToRelFsPath(path))
        // Converts the full path to the path relative to the current working directory.
        relPath, _ := filepath.Rel(cwd, fullPath)
        fmt.Println(relPath)
      }
    }
  }
  return nil
}

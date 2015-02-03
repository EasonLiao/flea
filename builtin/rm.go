package builtin

import (
  "flag"
  "fmt"
  "github.com/easonliao/flea/core"
  "os"
  "path/filepath"
)

func UsageRm() {
  usage :=
  `Usage: flea rm [--cached] (<file>|<directory>)

  --cached: use this option to unstage and remove paths only from the index.
            Working tree files, whether modified or not, will be left alone.
  `
  fmt.Println(usage)
  os.Exit(1)
}

func CmdRm() error {
  if len(os.Args) < 3 {
    // At least: flea rm <file>
    fmt.Println("Not enough arguments.")
    UsageRm()
  }
  flags := flag.NewFlagSet("rm", 0)
  cached := flags.Bool("cached", false, "delete file/directory from index tree")
  flags.Parse(os.Args[2:])

  rmPath := os.Args[len(os.Args) - 1]
  treePath := filepath.ToSlash(filepath.Join(core.GetPathPrefix(), rmPath))
  indexTree := core.GetIndexTree()

  if err := indexTree.Delete(treePath); err == nil {
    if *cached == false {
      // We need also deleting the path from working directory.
      fullPath := filepath.Join(core.GetRepoDirectory(), TreePathToRelFsPath(treePath))
      os.Remove(fullPath)
    }
  } else if err == core.ErrPathNotExist {
    fmt.Println("Can't find the path in index tree.")
    os.Exit(1)
  } else {
    fmt.Println(err.Error())
    os.Exit(1)
  }
  return nil
}

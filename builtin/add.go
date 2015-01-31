package builtin

import (
  "github.com/easonliao/flea/core"
  "os"
  "path"
  "path/filepath"
  "strings"
)

func CmdAdd() error {
  addPath := os.Args[len(os.Args) - 1]
  treePath := filepath.ToSlash(filepath.Join(core.GetPathPrefix(), addPath))
  add(treePath)
  return nil
}

func add(treePath string) error {
  fstree := core.GetFsTree()
  indextree := core.GetIndexTree()
  node, err := fstree.Get(treePath)
  if err != nil {
    return err
  }
  if !node.IsDir() {
    if hash, err := addFileToStore(treePath); err != nil {
      return err
    } else {
      return indextree.MkFileAll(treePath, hash)
    }
  } else {
    nodePaths := make([]string, 0, 64)
    fn := func(treePath string, node core.Node) error {
      if strings.HasPrefix(path.Base(treePath), ".") {
        if node.IsDir() {
          return core.SkipDirNode
        } else {
          return nil
        }
      }
      if !node.IsDir() {
        nodePaths = append(nodePaths, treePath)
      }
      return nil
    }
    // Finds all the paths of file node under treePath.
    fstree.Traverse(fn, treePath)
    if len(nodePaths) == 0 {
      return ErrEmptyDir
    }
    for _, nodePath := range(nodePaths) {
      if hash, err := addFileToStore(nodePath); err != nil {
        return err
      } else {
        if err := indextree.MkFileAll(nodePath, hash); err != nil {
          return err
        }
      }
    }
    return nil
  }
}

func addFileToStore(treePath string) ([]byte, error) {
  tree := core.GetFsTree()
  if node, err := tree.Get(treePath); err == nil {
    if node.IsDir() {
      return nil, ErrNotFile
    }
    data, err := node.GetData()
    if err != nil {
      return nil, err
    }
    hash, err := core.GetCAStore().StoreBlob(data)
    return hash, err
  } else {
    return nil, err
  }
}

package core

import (
  "bytes"
  "errors"
)

var (
  ErrEmptyTree = errors.New("core: tree is empty")
  ErrFileNotInCaStore = errors.New("core: file is not in CAStore")
)

// Builds a CATree from the staging area.
func BuildCATreeFromIndexFile() (*CATree, error) {
  idxTree := GetIndexTree()
  caStore := GetCAStore()
  root, _ := idxTree.Get("/")
  var rootHash []byte

  // Checks whether the index tree is empty.
  if bytes.Compare(root.GetHashValue(), EmptyDirHash[:]) == 0 {
    return nil, ErrEmptyTree
  }

  // Traverse the index tree and stores all the dir nodes to CAStore, it will also verify
  // that file nodes have already existed in CAStore.
  storeFn := func(treePath string, node Node) error {
    if node.IsDir() {
      if treePath == "/" {
        // Remembers the hash of root node.
        rootHash = node.GetHashValue()
      }
      hash, err := caStore.StoreTree([]byte(GetDirString(node)))
      if err != nil {
        return err
      }
      if bytes.Compare(hash, node.GetHashValue()) != 0 {
        // The hash value returned by CAStore should match the hash value calculated
        // by index tree.
        panic("The hash value of dir node doens't match returned by ca store")
      }
    } else {
      // The node is file, verifies it's in CAStore.
      if !caStore.Exists(node.GetHashValue()) {
        return ErrFileNotInCaStore
      }
    }
    return nil
  }

  err := indexTree.Traverse(storeFn, "/")
  if err != nil {
    return nil, err
  }
  if rootHash == nil {
    panic("The hash of root node is nil!")
  }
  return GetCATree(rootHash), nil
}

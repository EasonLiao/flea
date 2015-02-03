package core

import (
  "encoding/hex"
  "log"
  "strings"
)

// The tree structure stored in CAStore.
type CATree struct {
  root *CANode
}

// The node of CATree.
type CANode struct {
  hash  []byte
  children map[string]Node
}

// Gets the CATree from the hash value of root node.
func GetCATree(rootHash []byte) *CATree {
  return &CATree{newCANode(rootHash)}
}

// See Tree interface.
func (tree *CATree) Get(treePath string) (Node, error) {
  if treePath == "/" {
    return tree.root, nil
  }
  if !strings.HasPrefix(treePath, "/") {
    return nil, ErrInvalidPath
  }
  paths := strings.Split(treePath[1:], "/")
  var node Node = tree.root
  for _, name := range(paths) {
    if !node.IsDir() {
      return nil, ErrNotDir
    }
    children := node.GetChildren()
    if child, ok := children[name]; !ok {
      return nil, ErrPathNotExist
    } else {
      node = child
    }
  }
  return node, nil
}

// See Tree interface.
func (tree *CATree) Traverse(fn VisitFn, root string) error {
  node, err := tree.Get(root)
  if err != nil {
    return err
  }
  if node, ok := node.(*CANode); ok {
    return recursiveTraverse(root, node, fn)
  } else {
    panic("bug?")
  }
}

// Gets the hash value of root node.
func (tree *CATree) GetHash() []byte {
  return tree.root.GetHashValue()
}

func newCANode(hash []byte) *CANode {
  return &CANode{hash : hash}
}

// See Node interface.
func (node *CANode) GetHashValue() []byte {
  return node.hash
}

// See Node interface.
func (node *CANode) IsDir() bool {
  fType, _, err := GetCAStore().Get(node.GetHashValue())
  if err != nil {
    log.Fatal(err.Error())
  }
  return fType == TreeType
}

// See Node interface.
func (node *CANode) String() string {
  return String(node)
}

// See Node interface.
func (node *CANode) GetData() ([]byte, error) {
  fType, data, err := GetCAStore().Get(node.GetHashValue())
  if err != nil {
    return nil, err
  }
  if fType != BlobType {
    return nil, ErrNotFile
  }
  return data, nil
}

// See Node interface.
func (node *CANode) GetChildren() map[string]Node {
  if node.children != nil {
    return node.children
  }
  fType, data, err := GetCAStore().Get(node.GetHashValue())
  if err != nil {
    log.Fatal(err.Error())
  }
  if fType !=  TreeType {
    log.Fatal(ErrNotFile.Error())
  }
  if len(data) == 0 {
    // It's possible the directory is empty.
    return make(map[string]Node)
  }
  children := make(map[string]Node)
  dirString := string(data)
  rows := strings.Split(dirString, "\n")
  for _, row  := range(rows) {
    entries := strings.Split(row, " ")
    if len(entries) != 3 {
      log.Fatal("The number of entries per row is incorrect")
    }
    name := entries[2]
    hash, err := hex.DecodeString(entries[1])
    if err != nil {
      log.Fatal(err)
    }
    children[name] = newCANode(hash)
  }
  // Caches the children.
  node.children = children
  return children
}

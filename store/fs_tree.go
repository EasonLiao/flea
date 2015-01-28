package store

import (
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
)

var _ = log.Println

// FsTree is a wrapper on MemTree, it represents the tree structure of current working
// directory. It simply converts the namespace of working directory to the MemTree
// strucutre and provides the interface of Tree for accesses.
type FsTree struct {
  basePath string
  cache map[string]*FsNode
}

// Constructs a FsTree with the given path directory.
func NewFsTree(wd string) *FsTree {
  tree := &FsTree{basePath : wd, cache : make(map[string]*FsNode)}
  return tree
}

func (ft *FsTree) Get(path string) (Node, error) {
  if node, ok := ft.cache[path]; ok {
    return node, nil
  }
  node := NewFsTreeNode(filepath.Join(ft.basePath, path))
  if !node.IsExist() {
    return nil, ErrPathNotExist
  }
  ft.cache[path] = node
  return node, nil
}

func (ft *FsTree) Traverse(fn VisitFn) error {
  walkFn := func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    relPath, _ := filepath.Rel(ft.basePath, path)
    if relPath == "." {
      relPath = ""
    }
    relPath = "/" + relPath
    node, gerr := ft.Get(relPath)
    if gerr != nil {
      panic("bug?")
    }
    err = fn(relPath, node)
    if err == SkipDirNode {
      err = filepath.SkipDir
    }
    return err
  }
  return filepath.Walk(ft.basePath, walkFn)
}

type FsNode struct {
  path string
  hash []byte
}

func NewFsTreeNode(path string) *FsNode {
  return &FsNode{path, nil}
}

func (n *FsNode) GetHashValue() []byte {
  if n.hash != nil {
    return n.hash
  }
  if n.IsDir() {
    // If the node is the directory, the hash vlaue is the hash of dir string.
    hash, _, _ := WrapData(TreeType, []byte(GetDirString(n)))
    n.hash = hash[:]
  } else {
    // If it's a file, the hash value is the hash value of the file.
    data, err := ioutil.ReadFile(n.path)
    if err != nil {
      panic("Error while reading file " + n.path)
    }
    hash, _, _ := WrapData(BlobType, data)
    n.hash = hash[:]
  }
  return n.hash
}

func (n *FsNode) IsDir() bool {
  if fi, err := os.Stat(n.path); err == nil {
    return fi.IsDir()
  }
  panic("File not exists")
}

func (n *FsNode) GetChildren() map[string]Node {
  children := make(map[string]Node)
  walkFn := func(path string, info os.FileInfo, err error) error {
    name, _ := filepath.Rel(n.path, path)
    if name == "." {
      return nil
    }
    children[name] = NewFsTreeNode(path)
    if info.IsDir() {
      return filepath.SkipDir
    }
    return nil
  }
  filepath.Walk(n.path, walkFn)
  return children
}

func (n *FsNode) IsExist() bool {
  if _, err := os.Stat(n.path); err == nil {
    return true
  }
  return false
}

func (n *FsNode) String() string {
  return String(n)
}

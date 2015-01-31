package core

import (
  "io/ioutil"
  "os"
  "path/filepath"
)

var fsTree *FsTree

// FsTree implements Tree interface. It represents the tree structure of current working
// directory. It's also read-only.
type FsTree struct {
  baseFsPath string
  cache map[string]*FsTreeNode
}

// Gets the singleton FsTree.
func GetFsTree() *FsTree {
  if fsTree == nil {
    fsTree = newFsTree(GetWorkingDirectory())
  }
  return fsTree
}

// Constructs a FsTree with the given path directory.
func newFsTree(fsPath string) *FsTree {
  tree := &FsTree{baseFsPath :fsPath, cache : make(map[string]*FsTreeNode)}
  return tree
}

// See Tree interface.
func (ft *FsTree) Get(treePath string) (Node, error) {
  if node, ok := ft.cache[treePath]; ok {
    return node, nil
  }
  node := newFsTreeNode(filepath.Join(ft.baseFsPath, filepath.FromSlash(treePath)), ft)
  if !node.IsExist() {
    return nil, ErrPathNotExist
  }
  ft.cache[treePath] = node
  return node, nil
}

// See Tree interface.
func (ft *FsTree) Traverse(fn VisitFn, root string) error {
  fsPath := filepath.FromSlash(root)
  walkFn := func(fsPath string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    relPath, err := filepath.Rel(ft.baseFsPath, fsPath)
    if err != nil {
      return err
    }
    treePath := filepath.ToSlash(relPath)
    if  treePath == "." {
      treePath = ""
    }
    treePath = "/" + treePath
    node, gerr := ft.Get(treePath)
    if gerr != nil {
      panic("bug?")
    }
    err = fn(treePath, node)
    if err == SkipDirNode {
      err = filepath.SkipDir
    }
    return err
  }
  return filepath.Walk(filepath.Join(ft.baseFsPath, fsPath), walkFn)
}

type FsTreeNode struct {
  fsPath string
  hash []byte
  tree *FsTree
  children map[string]Node
}

func newFsTreeNode(fsPath string, tree* FsTree) *FsTreeNode {
  return &FsTreeNode{fsPath : fsPath, tree : tree}
}

func (n *FsTreeNode) GetHashValue() []byte {
  if n.hash != nil {
    return n.hash
  }
  if n.IsDir() {
    // If the node is the directory, the hash vlaue is the hash of dir string.
    hash, _, _ := WrapData(TreeType, []byte(GetDirString(n)))
    n.hash = hash[:]
  } else {
    // If it's a file, the hash value is the hash value of the file.
    data, err := ioutil.ReadFile(n.fsPath)
    if err != nil {
      panic("Error while reading file " + n.fsPath)
    }
    hash, _, _ := WrapData(BlobType, data)
    n.hash = hash[:]
  }
  return n.hash
}

func (n *FsTreeNode) IsDir() bool {
  if fi, err := os.Stat(n.fsPath); err == nil {
    return fi.IsDir()
  }
  panic("File not exists")
}

func (n *FsTreeNode) GetChildren() map[string]Node {
  if n.children != nil {
    return n.children
  }
  children := make(map[string]Node)
  walkFn := func(fsPath string, info os.FileInfo, err error) error {
    name, _ := filepath.Rel(n.fsPath, fsPath)
    if name == "." {
      return nil
    }
    relPath, err := filepath.Rel(n.tree.baseFsPath, fsPath)
    if err != nil {
      return err
    }
    treePath := filepath.ToSlash(relPath)
    children[name], err = n.tree.Get(treePath)
    if err != nil {
      panic(err.Error() + fsPath)
    }
    if info.IsDir() {
      return filepath.SkipDir
    }
    return nil
  }
  filepath.Walk(n.fsPath, walkFn)
  // Caches the children.
  n.children = children
  return children
}

func (n *FsTreeNode) GetData() ([]byte, error) {
  if n.IsDir() {
    return nil, ErrNotFile
  }
  data, err := ioutil.ReadFile(n.fsPath)
  return data, err
}

func (n *FsTreeNode) IsExist() bool {
  if _, err := os.Stat(n.fsPath); err == nil {
    return true
  }
  return false
}

func (n *FsTreeNode) String() string {
  return String(n)
}

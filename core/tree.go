package core

import (
  "bytes"
  "encoding/hex"
  "errors"
  "fmt"
  "log"
  "sort"
)

var (
  ErrNotDir = errors.New("store: not directory node")
  ErrPathNotExist = errors.New("store: path not exist")
  ErrInvalidPath = errors.New("store: invalid tree path")
  ErrNodeAlreadyExist = errors.New("store: file already exists")
  ErrReadOnlyRoot = errors.New("store: root node is read-only")
)

// The hash value of am empty directory.
var EmptyDirHash, _, _ = WrapData(TreeType, []byte(""))
// Returns SkipDirNode in VisitFn to skip traversing into directories.
var SkipDirNode = errors.New("store: skip traversing into directory")
// VisitFn defines the signature of function which will be invoked during tree traversal.
type VisitFn func(path string, node Node) error

// Tree interface.
type Tree interface {
  // Gets the node of the given path.
  Get(path string) (Node, error)

  // Traverses the tree structure.  VisitFn will be invoked during traversal.
  Traverse(fn VisitFn) error
}

// Node interface.
type Node interface {
  // Gets the hash value of the node.
  GetHashValue() []byte

  // Gets the children of dir node.
  GetChildren() map[string]Node

  // Checks if the node is directory.
  IsDir() bool

  // Converts node to readable string.
  String() string
}

// Compares two trees and returns the differences. bMisses is a list path of files which
// are included in tree a but not tree b, aMisses is a list path of files which are
// included tree b but not tree a, diffs is a list of path of files which are included
// in both trees but with different hash values.
func CompareTrees(a Tree, b Tree) (bMisses []string, aMisses []string, diffes []string) {
  misses := make([]string, 0, 64)
  diffes = make([]string, 0, 64)
  peerTree := b

  visitFn := func(path string, node Node) error {
    peerNode, err := peerTree.Get(path)
    if err == ErrPathNotExist {
      misses = append(misses, path)
      if node.IsDir() {
        return SkipDirNode
      }
      return nil
    }
    if node.IsDir() != peerNode.IsDir() {
      // One is directory, one is file.
      diffes = append(diffes, path)
      return SkipDirNode
    }
    // Two nodes have the same type.
    isHashSame := bytes.Compare(node.GetHashValue(), peerNode.GetHashValue()) == 0
    if isHashSame && node.IsDir() {
      // If two directories have the same hash value we don't need to traverse into the
      // directories.
      return SkipDirNode
    }
    if !node.IsDir() && !isHashSame {
      // They are two files with different hash values.
      diffes = append(diffes, path)
    }
    return nil
  }

  // Starts first traversal on tree a.
  a.Traverse(visitFn)
  bMisses = misses

  firstDiffes := diffes
  // Resets the closures so next time we'll start the traversal on tree b.
  misses = make([]string, 0, 64)
  diffes = make([]string, 0, len(firstDiffes))
  peerTree = a
  b.Traverse(visitFn)
  aMisses = misses
  if len(diffes) != len(firstDiffes) {
    panic("a bug?")
  }
  return
}

func GetDirString(node Node) string {
  var content string
  if node.IsDir() {
    children := node.GetChildren()
    names := make([]string, len(children))
    i := 0
    for k, _ := range(children) {
      names[i] = k
      i++
    }
    sort.Strings(names)
    for _, name := range(names) {
      child, ok := children[name]
      if !ok {
        panic("bug?")
      }
      if child.IsDir() {
        content += "tree "
      } else {
        content += "blob "
      }
      content += hex.EncodeToString(child.GetHashValue())
      content += " " + name + "\n"
    }
  } else {
    log.Fatal("File node doesn't contain any data in MemTree")
  }
  return content
}

func String(node Node) string {
  if node.IsDir() {
    return fmt.Sprintf("[Type:Dir, Hash:%x, Children: %d]", node.GetHashValue(), len(node.GetChildren()))
  } else {
    return fmt.Sprintf("[Type:File, Hash:%x]", node.GetHashValue())
  }
}

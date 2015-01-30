package core

import (
  "path/filepath"
  "testing"
)

func TestIndexTree(t *testing.T) {
  dir, _ := mkDir("index_tree_test")
  file := filepath.Join(dir, "index")
  tree, _ := newIndexTree(file)
  tree.MkDir("/foo1")
  tree.MkDir("/foo2")
  tree.MkFile("/foo1/file1", generateRandomHash())
  tree.MkFile("/foo1/file2", generateRandomHash())
  if _, err := tree.Get("/foo1"); err != nil {
    t.Error("/foo1 doesnt exist.")
  }
  if _, err := tree.Get("/foo2"); err != nil {
    t.Error("/foo2 doesnt exist.")
  }
  if _, err := tree.Get("/foo1/file1"); err != nil {
    t.Error("/foo1/file1 doesnt exist.")
  }
  if _, err := tree.Get("/foo1/file2"); err != nil {
    t.Error("/foo1/file2 doesnt exist.")
  }
  // Creates second tree, it will restore itself from the file.
  tree2, _ := newIndexTree(file)

  m1, m2, diffes := CompareTrees(tree, tree2)
  if len(m1) != 0 || len(m2) != 0 || len(diffes) != 0 {
    t.Error("Inconsistency between two trees.")
  }
}

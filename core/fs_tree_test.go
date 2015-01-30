package core

import (
  "testing"
)

var files = map[string][]byte {
  "/README" : []byte("this is read me"),
  "/src" : nil,
  "/src/hello.py" : []byte("print \"hello world\""),
}

func TestFsTree(t *testing.T) {
  dir, err := mkDir("test_fs_tree")
  if err != nil {
    panic(err.Error())
  }
  err = createTempFiles(dir, files)
  if err != nil {
    panic(err.Error())
  }
  tree := newFsTree(dir)
  nodes := make([]string, 0)
  fn := func(path string, node Node) error {
    if path == "/" {
      return nil
    }
    nodes = append(nodes, path)
    return nil
  }
  tree.Traverse(fn)
  if len(nodes) != len(files) {
    t.Error("Number of files/directories is incorrect")
  }
  for _, name := range(nodes) {
    if _, ok := files[name]; !ok {
      t.Errorf("File/Dir %s is not in dir %s\n", name, dir)
    }
  }

  // Creates a MemTree by iterating the FsTree and verify they are the same.
  memTree := NewMemTree()
  fn = func(path string, node Node) error {
    if node.IsDir() {
      memTree.MkDir(path)
    } else {
      memTree.MkFile(path, node.GetHashValue())
    }
    return nil
  }
  tree.Traverse(fn)
  m1, m2, diffes := CompareTrees(tree, memTree)

  if len(m1) != 0 || len(m2) != 0 || len(diffes) != 0 {
    t.Error("Inconsistency between two trees.")
  }
}

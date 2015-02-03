package core

import (
  "encoding/hex"
  "math/rand"
  "testing"
)

func generateRandomHash() []byte {
  var arr [HashSize]byte
  for i := 0; i < HashSize; i++ {
    arr[i] = byte(rand.Int())
  }
  return arr[:]
}

func TestMemTree(t *testing.T) {
  tree := NewMemTree()
  root, err := tree.Get("/")
  if err != nil {
    t.Error("Failed to get root node")
  }
  // Verifies the root node is dir.
  if !root.IsDir() {
    t.Error("Root node is not directory node")
  }
  // The root node shouldn't have children initially.
  if len(root.GetChildren()) != 0 {
    t.Error("The initial size of the chilren of root is not 0")
  }
  // The initial hash of root node.
  hash1 := hex.EncodeToString(root.GetHashValue())
  err = tree.MkFile("/foo", generateRandomHash())
  if err != nil {
    t.Error("Failed to create /foo")
  }
  if len(root.GetChildren()) != 1 {
    t.Error("The size of children of root node should be 1")
  }
  hash2 := hex.EncodeToString(root.GetHashValue())
  // Hash values of root node should be changed.
  if hash1 == hash2 {
    t.Error("Hash values of root node were not changed after creating a file")
  }

  err = tree.MkFile("/foo/bar", generateRandomHash())
  // Creating file should be failed since /foo is not a directory
  if err != ErrNotDir {
    t.Error("Expecting failure of creating /foo/bar.")
  }

  // Deletes /foo
  tree.Delete("/foo")
  hash3 := hex.EncodeToString(root.GetHashValue())
  // Since we restore the state a tree to initial state, the hash value of root node
  // should be restored.
  if hash1 != hash3 {
    t.Error("The hash values of two same trees don't match")
  }

  // The /foo shouldn't exist anymore.
  _, err = tree.Get("/foo")
  if err != ErrPathNotExist {
    t.Error("Expecting /foo not existing in tree")
  }

  tree.MkDir("/foo")
  tree.MkFile("/foo/bar", generateRandomHash())
  _, err = tree.Get("/foo")
  if err != nil {
   t.Error("/foo should exist in tree")
  }
  _, err = tree.Get("/foo/bar")
  if err != nil {
   t.Error("/foo/bar should exist in tree")
  }

  // Deletes dir /foo, this should also delete /foo/bar
  tree.Delete("/foo")
  hash4 := hex.EncodeToString(root.GetHashValue())
  // Since we restore the state a tree to initial state, the hash value of root node
  // should be restored.
  if hash1 != hash4 {
    t.Error("The hash values of two same trees don't match")
  }

  // Can't delete root node.
  err = tree.Delete("/")
  if err != ErrReadOnlyRoot {
    t.Error("Expecting read-only root node error")
  }

  // Tests creating dir/file recursively.
  treeA := NewMemTree()
  treeB := NewMemTree()
  treeA.MkDirAll("/d1/d2/d3")
  treeB.MkDir("/d1")
  treeB.MkDir("/d1/d2")
  treeB.MkDir("/d1/d2/d3")
  m1, m2, diffes := CompareTrees(treeA, treeB)
  if len(m1) != 0 || len(m2) != 0 || len(diffes) != 0 {
    t.Error("Inconsistency between two trees.")
  }
  hash := generateRandomHash()
  treeA.MkFileAll("/d/dd/ddd", hash)
  treeB.MkDir("/d")
  treeB.MkDir("/d/dd")
  treeB.MkFile("/d/dd/ddd", hash)
  m1, m2, diffes = CompareTrees(treeA, treeB)
  if len(m1) != 0 || len(m2) != 0 || len(diffes) != 0 {
    t.Error("Inconsistency between two trees.")
  }

  // Tests clear.
  treeA.Clear()
  node, _ := treeA.Get("/")
  if len(node.GetChildren()) != 0 {
    t.Error("Error in clear")
  }
}

func TestTraversal(t *testing.T) {
  tree := NewMemTree()
  tree.MkFile("/foo", generateRandomHash())
  tree.MkDir("/fooDir")
  tree.MkFile("/fooDir/bar", generateRandomHash())

  visitMap := make(map[string]bool)
  visitFn := func(path string, node Node) error {
    visitMap[path] = true
    return nil
  }
  tree.Traverse(visitFn, "/")
  if _, ok := visitMap["/"]; !ok {
    t.Error("/ should be visited")
  }
  if _, ok := visitMap["/foo"]; !ok {
    t.Error("/foo should be visited")
  }
  if _, ok := visitMap["/fooDir"]; !ok {
    t.Error("/fooDir should be visited")
  }
  if _, ok := visitMap["/fooDir/bar"]; !ok {
    t.Error("/fooDir/bar should be visited")
  }
  if len(visitMap) != 4 {
    t.Error("The number of visited nodes is incorrect")
  }

  // Tests skipping traversing into /fooDir
  visitMap = make(map[string]bool)
  visitFn = func(path string, node Node) error {
    visitMap[path] = true
    if path == "/fooDir" {
      return SkipDirNode
    }
    return nil
  }
  tree.Traverse(visitFn, "/")
  if _, ok := visitMap["/"]; !ok {
    t.Error("/ should be visited")
  }
  if _, ok := visitMap["/foo"]; !ok {
    t.Error("/foo should be visited")
  }
  if _, ok := visitMap["/fooDir"]; !ok {
    t.Error("/fooDir should be visited")
  }
  if _, ok := visitMap["/fooDir/bar"]; ok {
    t.Error("/fooDir/bar should not be visited")
  }
  if len(visitMap) != 3 {
    t.Error("The number of visited nodes is incorrect")
  }
}

func TestCompareTrees1(t *testing.T) {
  treeA := NewMemTree()
  treeB := NewMemTree()

  hash := generateRandomHash()
  // Creates /foo file with same value on both trees.
  treeA.MkFile("/foo", hash)
  treeB.MkFile("/foo", hash)

  // Creates /dir directory on both trees.
  treeA.MkDir("/dir")
  treeB.MkDir("/dir")
  // Creates /dir2 directory on both trees.
  treeA.MkDir("/dir2")
  treeB.MkDir("/dir2")

  hash = generateRandomHash()
  // Creates /dir/foo file with the same hash value on both trees.
  treeA.MkFile("/dir/foo", hash)
  treeB.MkFile("/dir/foo", hash)

  hash = generateRandomHash()
  // Creates /dir/foo file with the same hash value on both trees.
  treeA.MkFile("/dir2/foo", hash)
  treeB.MkFile("/dir2/foo", hash)

  // Creates /dir/barA on treeA.
  treeA.MkFile("/dir/barA", generateRandomHash())
  // Creates /dir/barB on treeB.
  treeB.MkFile("/dir/barB", generateRandomHash())

  // Creates /dir/bar file with differen hash values on both trees.
  treeA.MkFile("/dir/bar", generateRandomHash())
  treeB.MkFile("/dir/bar", generateRandomHash())
  bMisses, aMisses, diffes := CompareTrees(treeA, treeB)

  if len(bMisses) != 1 {
    t.Error("The number of missed files on tree B is incorrect.")
  }
  if len(aMisses) != 1 {
    t.Error("The number of missed files on tree A is incorrect.")
  }
  if len(diffes) != 1 {
    t.Error("The number of different files between A and B is incorrect.")
  }
  if bMisses[0] != "/dir/barA" {
    t.Error("B should miss file /dir/barA")
  }
  if aMisses[0] != "/dir/barB" {
    t.Error("A should miss file /dir/barB")
  }
  if diffes[0] != "/dir/bar" {
    t.Error("the different file between A and B should be /dir/bar")
  }
}

func TestCompareTrees2(t *testing.T) {
  treeA := NewMemTree()
  treeB := NewMemTree()

  hash := generateRandomHash()
  // Creates /foo file with same value on both trees.
  treeA.MkFile("/foo", hash)
  treeB.MkFile("/foo", hash)

  treeA.MkDir("/dirA")
  treeA.MkFile("/dirA/foo1", generateRandomHash())
  treeA.MkFile("/dirA/foo2", generateRandomHash())

  treeB.MkDir("/dirB")
  treeB.MkFile("/dirB/foo1", generateRandomHash())
  treeB.MkFile("/dirB/foo2", generateRandomHash())

  bMisses, aMisses, diffes := CompareTrees(treeA, treeB)
  if len(bMisses) != 1 && len(aMisses) != 1 && len(diffes) != 0 {
    t.Error("Incorrect test result")
  }
  if bMisses[0] != "/dirA" && aMisses[0] != "/dirB" {
    t.Error("Incorrect test result")
  }
}

func TestSerialization(t *testing.T) {
  tree := NewMemTree()
  tree.MkFile("/foo",  generateRandomHash())
  tree.MkDir("/dir")
  tree.MkFile("/dir/foo1", generateRandomHash())
  tree.MkFile("/dir/foo2", generateRandomHash())

  data, err := tree.Serialize()
  if err != nil {
    t.Error("Error found in serialization of MemTree.")
  }
  newtree, err := Deserialize(data)
  if err != nil {
    t.Error("Error found in deserialization of MemTree.")
  }
  m1, m2, diffes := CompareTrees(tree, newtree)
  if len(m1) != 0 || len(m2) != 0 || len(diffes) != 0 {
    t.Error("Inconsistency after serialization/deserialization.")
  }
}

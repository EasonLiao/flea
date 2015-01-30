package core

import (
  "encoding/json"
  "path"
  "strings"
)

// Function signature of parameter of apply function.
type Op func(node *MemTreeNode) (changed bool, ret interface{}, err error)

// Represents the file system as a standard tree structure in memory.
// It doesn't contain the data of files, it only contains the namespace structure
// of the file system and the hash values of the actual data and directory.
// NOTE : The MemTree is not thread-safe.
type MemTree struct  {
  root *MemTreeNode
}

// Gets the node for a given path.
func (mt *MemTree) Get(treePath string) (Node, error) {
  if treePath == "/" {
    return mt.root, nil
  }
  nodeName := path.Base(treePath)
  dir := path.Dir(treePath)
  op := func(node *MemTreeNode) (changed bool, ret interface{}, err error) {
    if !node.Dir {
      err = ErrNotDir
      return
    }
    childNode, ok := node.Children[nodeName]
    if !ok {
      err = ErrPathNotExist
      return
    }
    ret = childNode
    // Get op is read-only.
    changed = false
    return
  }
  ret, err := mt.apply(dir, op)
  if err != nil {
    return nil, err
  }
  if node, ok := ret.(*MemTreeNode); !ok {
    panic("bug?")
  } else {
    return node, nil
  }
}

// Traverse the tree structure. MemTree traverses the tree in DFS way.
func (mt *MemTree) Traverse(fn VisitFn) error {
  return recursiveTraverse("/", mt.root, fn)
}

// Creates a directory in tree.
func (mt *MemTree) MkDir(treePath string) (err error) {
  if treePath == "/" {
    err = ErrReadOnlyRoot
    return
  }
  dirName := path.Base(treePath)
  dir := path.Dir(treePath)
  op := func(node *MemTreeNode) (changed bool, ret interface{}, err error) {
    if !node.Dir {
      err = ErrNotDir
      return
    }
    _, ok := node.Children[dirName]
    if ok {
      err = ErrNodeAlreadyExist
      return
    }
    node.Children[dirName] = newDirMemTreeNode()
    changed = true
    return
  }
  _, err = mt.apply(dir, op)
  return
}

// Creates a file with given hash value in tree.
func (mt *MemTree) MkFile(treePath string, hash []byte) (err error) {
  if treePath == "/" {
    err = ErrReadOnlyRoot
    return
  }
  fileName := path.Base(treePath)
  dir := path.Dir(treePath)
  op := func(node *MemTreeNode) (changed bool, ret interface{}, err error) {
    if !node.Dir {
      err = ErrNotDir
      return
    }
    _, ok := node.Children[fileName]
    if ok {
      err = ErrNodeAlreadyExist
      return
    }
    node.Children[fileName] = newFileMemTreeNode(hash)
    changed = true
    return
  }
  _, err = mt.apply(dir, op)
  return
}

// Deletes a node from the tree. If the node is a directory the whole directory will be
// deleted.
func (mt *MemTree) Delete(treePath string) (err error) {
  if treePath == "/" {
    err = ErrReadOnlyRoot
    return
  }
  nodeName := path.Base(treePath)
  dir := path.Dir(treePath)
  op := func(node *MemTreeNode) (changed bool, ret interface{}, err error) {
    if !node.Dir {
      err = ErrNotDir
      return
    }
    _, ok := node.Children[nodeName]
    if !ok {
      err = ErrPathNotExist
      return
    }
    delete(node.Children, nodeName)
    changed = true
    return
  }
  _, err = mt.apply(dir, op)
  return
}

// Serializes the MemTree to byte array.
func (mt *MemTree) Serialize() ([]byte, error) {
  type tuple struct {
    Path string `json:path`
    Hash []byte `json:hash`
  }
  nodes := make([]tuple, 0)
  traverseFn := func(treePath string, node Node) error {
    if node.IsDir() {
      nodes = append(nodes, tuple{treePath, nil})
    } else {
      nodes = append(nodes, tuple{treePath, node.GetHashValue()})
    }
    return nil
  }
  mt.Traverse(traverseFn)
  data, err := json.Marshal(nodes)
  if err != nil {
    return nil, err
  }
  return data, nil
}

// Creates a MemTree.
func NewMemTree() *MemTree {
  return &MemTree{newDirMemTreeNode()}
}

// Deserializes the byte array to MemTree.
func Deserialize(data []byte) (*MemTree, error) {
  type tuple struct {
    Path string `json:path`
    Hash []byte `json:hash`
  }
  nodes := make([]tuple, 0)
  err := json.Unmarshal(data, &nodes)
  if err != nil {
    return nil, err
  }
  tree := NewMemTree()
  for _, t := range(nodes) {
    if t.Hash != nil {
      tree.MkFile(t.Path, t.Hash)
    } else {
      tree.MkDir(t.Path)
    }
  }
  return tree, nil
}

// Apply an operation to node of the given path. This is the primitive used by other
// methods like Delete, MkDir, MkFile. This separates the traversing of the tree from
// the actual operation on the tree.
func (mt *MemTree) apply(treePath string, op Op) (ret interface{}, err error) {
  if !strings.HasPrefix(treePath, "/") {
    // path must starts with /
    err = ErrInvalidPath
    return
  }
  // Trim the root path
  treePath= treePath[1:]
  _, ret, err = recursive(mt.root,treePath, op)
  return
}

// Recursive traverse. Used by apply method only.
func recursive(node *MemTreeNode, remPath string, op Op) (changed bool, ret interface{}, err error) {
  if remPath == "" {
    // Last node in remPath, invokes op function.
    changed, ret, err = op(node)
  } else {
    // We need go into its children.
    if !node.Dir {
      err = ErrNotDir
      return
    }
    var childName string
    sep := strings.Index(remPath, "/")
    if sep != -1 {
      childName, remPath = remPath[:sep], remPath[sep+1:]
    } else {
      childName = remPath
      remPath = ""
    }
    childNode, ok := node.Children[childName]
    if !ok {
      err = ErrPathNotExist
      return
    }
    changed, ret, err = recursive(childNode, remPath, op)
  }
  if err != nil {
    return
  }
  if changed {
    node.updateHashValue()
  }
  return
}

func recursiveTraverse(treePath string, node* MemTreeNode, fn VisitFn) (err error) {
  err = fn(treePath, node)
  if err == nil && node.Dir {
    // node is directory, traverse into it.
    for childName, childNode := range(node.Children) {
      err = recursiveTraverse(path.Join(treePath, childName), childNode, fn)
      if err != nil {
        return err
      }
    }
  }
  if err == SkipDirNode {
    // SkipDirNode is not actual error, reset it to nil.
    err = nil
  }
  return
}

// Node of MemTree.
type MemTreeNode struct {
  //  Whether the node is directory or not.
  Dir bool
  // Hash value of the node.
  Hash [HashSize]byte
  // The children of the node if it's the directory.
  Children map[string]*MemTreeNode
}

func newDirMemTreeNode() *MemTreeNode {
  return &MemTreeNode{true, EmptyDirHash, make(map[string]*MemTreeNode)}
}

func newFileMemTreeNode(hash []byte) *MemTreeNode {
  if len(hash) != HashSize {
    panic("invalid length of hash")
  }
  var hashArr  [HashSize]byte
  copy(hashArr[:], hash)
  return &MemTreeNode{false, hashArr, nil}
}

func (n *MemTreeNode) GetHashValue() []byte {
  return n.Hash[:]
}

func (n *MemTreeNode) GetChildren() map[string]Node {
  children := make(map[string]Node)
  for k, node := range(n.Children) {
    children[k] = node
  }
  return children
}

func (n *MemTreeNode) IsDir() bool {
  return n.Dir
}

func (n *MemTreeNode) updateHashValue() {
  if n.Dir {
    hash, _, _ := WrapData(TreeType, []byte(GetDirString(n)))
    n.Hash = hash
  }
}

func (n *MemTreeNode) GetData() ([]byte, error) {
  panic("MemTreeNode doesn't support GetData() method.")
}

// Convert the Node to string.
func (n *MemTreeNode) String() string {
  return String(n)
}

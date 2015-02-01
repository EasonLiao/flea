package core

import (
  "bytes"
  "encoding/hex"
  "encoding/json"
  "errors"
  "io/ioutil"
  "log"
  "path/filepath"
)

var (
  ErrEmptyTree = errors.New("core: tree is empty")
  ErrFileNotInCaStore = errors.New("core: file is not in CAStore")
)

type Commit struct {
  Tree        []byte  `json:tree`
  PrevCommit  []byte  `json:prev`
  Author      string  `json:author`
  Comment     string  `json:comment`
}

func (c *Commit) GetPrevCommit() *Commit {
  if c.PrevCommit == nil {
    // No ancester commit object.
    return nil
  }
  fType, data, err := GetCAStore().Get(c.PrevCommit)
  if err != nil {
    panic(err.Error())
  }
  if fType != CommitType {
    panic("Not a valid commit hash.")
  }
  commit, err := fromBytesToCommitObject(data)
  if err != nil {
    panic(err.Error())
  }
  return commit
}

func (c* Commit) GetCommitHash() []byte {
  data, err := fromCommitObjectToBytes(c)
  if err != nil {
    panic("Not valid commit object")
  }
  hash, _, _ := WrapData(CommitType, data)
  return hash[:]
}

// Creates a commit object in CAStore.
func CreateCommitObject(tree, prevCommit []byte, author, comment string) ([]byte, error) {
  commit := Commit{tree, prevCommit, author, comment}
  store := GetCAStore()
  if !store.Exists(tree) {
    log.Fatal("Invalid tree hash.")
  }
  if prevCommit != nil && !store.Exists(prevCommit) {
    log.Fatal("Invalid previous commit  hash.")
  }
  data, err := fromCommitObjectToBytes(&commit)
  if err != nil {
    return nil, err
  }
  return GetCAStore().StoreCommit(data)
}

// Updates the head commit of a the branch.
func UpdateBranchHead(branch string, commitHash []byte) {
  if !GetCAStore().Exists(commitHash) {
    panic("Not a valid commit hash.")
  }
  hashString := hex.EncodeToString(commitHash)
  if err := ioutil.WriteFile(filepath.Join(GetBranchHeadDir(), branch), []byte(hashString), 0777); err != nil {
    panic("Failed to update the head of the branch." + err.Error())
  }
}

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

func fromCommitObjectToBytes(commit *Commit) ([]byte, error) {
  return json.Marshal(commit)
}

func fromBytesToCommitObject(data []byte) (*Commit, error) {
  var commit Commit
  err := json.Unmarshal(data, &commit)
  return &commit, err
}

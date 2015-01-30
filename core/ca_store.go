// Core implementations of Flea.
package core

import (
  "bytes"
  "crypto/sha1"
  "encoding/hex"
  "errors"
  "io/ioutil"
  "os"
  "path"
  "path/filepath"
  "strconv"
  "strings"
)

var (
  ErrNoMatch = errors.New("store: no match files")
  ErrHashTooShort = errors.New("store: hash value is too short(more than one files match)")
  ErrNotValidHash = errors.New("store: not valid hash value")
  ErrFileCorrupted = errors.New("store: file corrupted")
  ErrInvalidType = errors.New("store: invalid file type")
)

const (
  BlobType = "blob"
  TreeType = "tree"
  CommitType = "commit"
)

// The length of hash value.
const HashSize = 20
// CAStore is singleton object.
var caStore *CAStore = nil

// Content-addressable store.
type CAStore struct {
  dir string
}

// Gets the CAStore instance, it's a singleton object.
func GetCAStore() *CAStore {
  if caStore == nil {
    caStore = newCAStore(GetStoreDirectory())
  }
  return caStore
}

// Stores blob data to content-addressable store.
func (store *CAStore) StoreBlob(data []byte) ([]byte, error) {
  hash, blob, err := WrapData(BlobType, data)
  if err != nil {
    return nil, err
  }
  fileName := hex.EncodeToString(hash[:])
  ioutil.WriteFile(filepath.Join(store.dir, fileName), blob, 0444)
  return hash[:], nil
}

// Gets a list of filenames that match the prefix of the hash value.
func (store *CAStore) GetFileName(hashPrefix []byte) (names []string, err error) {
  if len(hashPrefix) > HashSize {
    err = ErrNotValidHash
    return
  }
  names = make([]string, 0, 1)
  hashString := hex.EncodeToString(hashPrefix)
  walkFun := func(path string, info os.FileInfo, err error) error {
    if info.IsDir() && path != store.dir {
      return filepath.SkipDir
    }
    name := info.Name()
    if strings.HasPrefix(name, hashString) {
      names = append(names, name)
    }
    return nil
  }
  filepath.Walk(store.dir, filepath.WalkFunc(walkFun))
  return
}

// Returns the type and data of the file based on the hash prefix.
func (store *CAStore) Get(hashPrefix []byte) (name string, fileType string, data []byte, err error) {
  names, err := store.GetFileName(hashPrefix)
  if err != nil {
    return
  }
  if len(names) > 1 {
    err = ErrHashTooShort
    return
  } else if len(names) == 0 {
    err = ErrNoMatch
    return
  }
  name = names[0]
  data, err = ioutil.ReadFile(path.Join(store.dir, name))
  sepIdx := bytes.IndexByte(data, 0)
  header, data := data[:sepIdx], data[sepIdx + 1:]
  headers := strings.Split(string(header), " ")
  fileType = headers[0]
  length, err := strconv.Atoi(headers[1])
  if err != nil {
    return
  }
  // sanity check, length field must match the actual length of data.
  if length != len(data) {
    err = ErrFileCorrupted
  }
  return
}

func WrapData(fileType string, data []byte) (hash [HashSize]byte, blob []byte, err error) {
  if fileType != BlobType && fileType != TreeType && fileType != CommitType {
    err =  ErrInvalidType
    return
  }
  var buffer bytes.Buffer
  buffer.WriteString(fileType)
  buffer.WriteString(" ")
  buffer.WriteString(strconv.Itoa(len(data)))
  buffer.WriteByte(0)
  buffer.Write(data)
  blob = buffer.Bytes()
  hash = sha1.Sum(blob)
  return
}

func newCAStore(dir string) *CAStore {
  return &CAStore{dir : dir}
}

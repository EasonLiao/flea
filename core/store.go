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

var storeDir string
var isDirInit bool = false

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

const hashSize = 20

func InitStoreDir(dir string) {
  storeDir = dir
  isDirInit = true
}

func GetStoreDir() string {
  return storeDir
}

func WrapData(fileType string, data []byte) (hash [hashSize]byte, blob []byte, err error) {
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

func StoreBlob(data []byte) [20]byte {
  assertDirInit()
  hash, blob, _ := WrapData(BlobType, data)
  fileName := hex.EncodeToString(hash[:])
  ioutil.WriteFile(filepath.Join(storeDir, fileName), blob, 0444)
  return hash
}

// Returns the type and data of the file based on the hash prefix.
func Get(hashPrefix []byte) (name string, fileType string, data []byte, err error) {
  names, err := GetFileName(hashPrefix)
  if err != nil {
    return
  }
  if len(names) > 1 {
    err = ErrHashTooShort
    return
  }
  name = names[0]
  data, err = ioutil.ReadFile(path.Join(storeDir, name))
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

// Given the prefix of the hash, returns a list of file names which have the
// prefix. Returns err if no file matches or the hash prefix is invalid.
func GetFileName(hashPrefix []byte) (names []string, err error) {
  assertDirInit()
  if len(hashPrefix) > hashSize {
    err = ErrNotValidHash
    return
  }
  names = make([]string, 0, 1)
  hashString := hex.EncodeToString(hashPrefix)
  walkFun := func(path string, info os.FileInfo, err error) error {
    if info.IsDir() && path != storeDir {
      return filepath.SkipDir
    }
    name := info.Name()
    if strings.HasPrefix(name, hashString) {
      names = append(names, name)
    }
    return nil
  }
  filepath.Walk(storeDir, filepath.WalkFunc(walkFun))
  if len(names) == 0 {
    err = ErrNoMatch
  }
  return
}

func assertDirInit() {
  if !isDirInit {
    panic("Dir has not been initialized.")
  }
}

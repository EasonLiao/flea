// Core implementations of Flea.
package core

import (
  "bytes"
  "crypto/sha1"
  "encoding/hex"
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

var (
  ErrNoMatch = errors.New("core: no match files")
  ErrHashTooShort = errors.New("core: hash value is too short(more than one files match)")
  ErrNotValidHash = errors.New("core: not valid hash value")
  ErrFileCorrupted = errors.New("core: file corrupted")
  ErrInvalidType = errors.New("core: invalid file type")
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
  store.write(fileName, blob)
  return hash[:], nil
}

// Stores tree data to content-addressable store.
func (store *CAStore) StoreTree(data []byte) ([]byte, error) {
  hash, blob, err := WrapData(TreeType, data)
  if err != nil {
    return nil, err
  }
  fileName := hex.EncodeToString(hash[:])
  store.write(fileName, blob)
  return hash[:], nil
}

// Stores commit data to content-addressable store.
func (store *CAStore) StoreCommit(data []byte) ([]byte, error) {
  hash, blob, err := WrapData(CommitType, data)
  if err != nil {
    return nil, err
  }
  fileName := hex.EncodeToString(hash[:])
  store.write(fileName, blob)
  return hash[:], nil
}

// Gets a list of full hashs that match the prefix of the hash value. The return values can be:
// 1) a list of hashs, nil
// 2) undefined, ErrNotValidHash
func (store *CAStore) GetMatchedHashs(hashPrefix []byte) (hashs [][]byte, err error) {
  if len(hashPrefix) > HashSize {
    err = ErrNotValidHash
    return
  }
  hashs = make([][]byte, 0, 1)
  hashString := hex.EncodeToString(hashPrefix)
  walkFun := func(path string, info os.FileInfo, err error) error {
    if info.IsDir() && path != store.dir {
      return filepath.SkipDir
    }
    name := info.Name()
    if strings.HasPrefix(name, hashString) {
      if hash, err := hex.DecodeString(name); err != nil {
      } else {
        hashs = append(hashs, hash)
      }
    }
    return nil
  }
  filepath.Walk(store.dir, filepath.WalkFunc(walkFun))
  return
}

// Given the hash value, gets the content stored in CAStore. The return values can be:
// 1) fileType, data, nil
// 2) "", nil, ErrNoMatch
func (store* CAStore) Get(hash []byte) (fileType string, data []byte, err error) {
  fileName := hex.EncodeToString(hash)
  fullPath := filepath.Join(store.dir, fileName)
  if data, err = ioutil.ReadFile(fullPath); err == nil {
    var header []byte
    sepIdx := bytes.IndexByte(data, 0)
    header, data = data[:sepIdx], data[sepIdx + 1:]
    headers := strings.Split(string(header), " ")
    fileType = headers[0]
    length, err := strconv.Atoi(headers[1])
    if err != nil {
      fmt.Println("Failed to conver %s to integer.", headers[1])
      os.Exit(1)
    }
    // Sanity check, length field must match the actual length of data.
    if length != len(data) {
      fmt.Println("The length is not correct, %s is invalid file", fileName)
      os.Exit(1)
    }
  } else if os.IsNotExist(err) {
    err = ErrNoMatch
  }
  return
}

// Given a hash value, returns true if a file with the given hash exists in store.
func (store *CAStore) Exists(hash []byte) bool {
  fileName := hex.EncodeToString(hash)
  if _, err := os.Stat(filepath.Join(store.dir, fileName)); err == nil {
    return true
  }
  return false
}

// Write data to CAStore. The fileName is just the hash string.
func (store *CAStore) write(fileName string, data []byte) {
  hash := sha1.Sum(data)
  if fileName != hex.EncodeToString(hash[:]) {
    // Sanity check, verifies the fileName is correct for the given data.
    fmt.Println("Hash of the data doesn't match the file name.")
    os.Exit(1)
  }
  fullPath := filepath.Join(store.dir, fileName)
  if _, err := os.Stat(fullPath); err == nil {
    // The file has alredy existed.
    return
  } else if os.IsNotExist(err) {
    if err := ioutil.WriteFile(filepath.Join(store.dir, fileName), data, 0444); err != nil {
      fmt.Println("Failed to write file.")
      os.Exit(1)
    }
  } else {
    fmt.Printf("Err: %s", err.Error())
    os.Exit(1)
  }
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

package store

import (
  "bytes"
  "encoding/hex"
  //"log"
  "os"
  "testing"
)

func TestStore(t *testing.T) {
  InitStoreDir(os.TempDir())
  data := []byte("what is up, doc?")
  hash1 := StoreBlob(data)
  hash2 := StoreBlob(data)
  // The hash of the two should be the same.
  if bytes.Compare(hash1[:], hash2[:]) != 0 {
    t.Error("Hash values don't match for same content")
  }
  name, fileType, content, err := Get(hash1[:])
  if err != nil {
    t.Error("error:", err.Error)
  }
  if name != hex.EncodeToString(hash1[:]) {
    t.Error("Names don't match")
  }
  if fileType != BlobType {
    t.Error("Invalid type")
  }
  if bytes.Compare(data, content) != 0 {
    t.Error("Data doesn't match")
  }
}

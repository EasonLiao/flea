package core

import (
  "bytes"
  "encoding/hex"
  "testing"
)

func TestStore(t *testing.T) {
  dir, _ := mkDir("ca_store")
  store := newCAStore(dir)
  data := []byte("what is up, doc?")
  hash1, _ := store.StoreBlob(data)
  hash2, _ := store.StoreBlob(data)
  // The hash of the two should be the same.
  if bytes.Compare(hash1[:], hash2[:]) != 0 {
    t.Error("Hash values don't match for same content")
  }
  name, fileType, content, err := store.GetWithPrefix(hash1[:])
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

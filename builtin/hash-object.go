package builtin

import (
  "github.com/easonliao/flea/store"
)

func CmdHashObject() int {
  store.StoreBlob([]byte("hello"))
  return 0
}

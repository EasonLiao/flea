package builtin

import (
  "flag"
  "fmt"
  "github.com/easonliao/flea/core"
  "io/ioutil"
  "log"
  "os"
)

func CmdHashObject() error {
  flags := flag.NewFlagSet("hash-object", 0)
  // The first one is the executable name, the second one is the command name,
  // actual arguments start from the third one.
  flags.Parse(os.Args[2:])
  data, err := ioutil.ReadAll(os.Stdin)
  // Trim the last character (EOF)?
  if err != nil {
    log.Fatal(err)
  }
  store := core.GetCAStore()
  hash, _ := store.StoreBlob(data)
  fmt.Printf("%x\n", hash)
  return nil
}

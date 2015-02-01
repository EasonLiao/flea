package builtin

import (
  "encoding/hex"
  "flag"
  "fmt"
  "github.com/easonliao/flea/core"
  "log"
  "os"
)

func CmdCatFile() error {
  if len(os.Args) <= 2 {
    log.Fatal("Not enough arguments.")
  }
  flags := flag.NewFlagSet("cat-file", 0)
  printType := flags.Bool("t", false, "file type")
  flags.Parse(os.Args[2:])
  hashPrefix :=  os.Args[len(os.Args) - 1]
  hash, err := hex.DecodeString(hashPrefix)
  if err != nil {
    log.Fatal("Invalid hash values.")
  }
  store := core.GetCAStore()
  _, fileType, data, err := store.GetWithPrefix(hash)
  if err != nil {
    log.Fatal(err)
  }
  if *printType {
    fmt.Println(fileType)
  } else {
    fmt.Println(string(data))
  }
  return nil
}

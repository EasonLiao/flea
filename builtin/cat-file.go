package builtin

import (
  "encoding/hex"
  "flag"
  "fmt"
  "github.com/easonliao/flea/store"
  "log"
  "os"
)

func CmdCatFile() int {
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
  _, fileType, data, err := store.Get(hash)
  if err != nil {
    log.Fatal(err)
  }
  if *printType {
    fmt.Println(fileType)
  } else {
    fmt.Println(string(data))
  }
  return 0
}

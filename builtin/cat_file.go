package builtin

import (
  "encoding/hex"
  "flag"
  "fmt"
  "github.com/easonliao/flea/core"
  "os"
)

func CmdCatFile() error {
  if len(os.Args) <= 2 {
    fmt.Println("Not enough arguments.")
    os.Exit(1)
  }
  flags := flag.NewFlagSet("cat-file", 0)
  printType := flags.Bool("t", false, "file type")
  flags.Parse(os.Args[2:])
  hashPrefix :=  os.Args[len(os.Args) - 1]
  hash, err := hex.DecodeString(hashPrefix)
  if err != nil {
    fmt.Println("Invalid hash values.")
    os.Exit(1)
  }
  store := core.GetCAStore()
  hashs := store.GetMatchedHashs(hash)
  if len(hashs) > 1 {
    fmt.Printf("More than one file match %s\n", hashPrefix)
    os.Exit(1)
  } else if len(hashs) == 1 {
    fType, data, err := store.Get(hashs[0])
    if err != nil {
      fmt.Printf("Error: %s\n", err.Error())
      os.Exit(1)
    }
    if *printType {
      fmt.Println(fType)
    } else {
      fmt.Println(string(data))
    }
  }
  return nil
}

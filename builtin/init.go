package builtin

import (
  "fmt"
  "github.com/easonliao/flea/core"
  "os"
)

func UsageInit() {
  fmt.Println("flea init")
  os.Exit(1)
}

func CmdInit() error {
  err := core.InitNew()
  return err
}

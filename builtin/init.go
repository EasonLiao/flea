package builtin

import (
  "github.com/easonliao/flea/core"
)

func CmdInit() error {
  err := core.InitNew()
  return err
}

package builtin

import (
  "github.com/easonliao/flea/core"
)

func CmdInit() int {
  core.InitNew()
  return 0
}

package main

import (
  "fmt"
  "github.com/easonliao/flea/builtin"
  "github.com/easonliao/flea/setup"
  "log"
  "os"
  "strings"
)

// debuggin only
var _ = fmt.Println

const (
  flagNeedSetup = (1 << iota)
  flagEnd
)

type cmdStruct struct {
  fun func() int
  flag int
}

var commandsTable = map[string]cmdStruct {
  "init" : {fun : builtin.CmdInit},
  "hash-object" : {fun : builtin.CmdHashObject, flag : flagNeedSetup},
  "cat-file" : {fun : builtin.CmdCatFile, flag : flagNeedSetup},
}

func runBuiltin(cmd string) {
  if cmdSt, ok := commandsTable[cmd]; !ok {
    log.Fatal("Unkown command")
  } else {
    if cmdSt.flag & flagNeedSetup != 0 {
      setup.SetupFleaDir()
    }
    cmdSt.fun()
  }
}

func main() {
  exe := os.Args[0]
  if exe != "flea" {
    panic("Invalid executable name.")
  }
  if (len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-")) {
    runBuiltin(os.Args[1])
  } else {
    log.Println("Else")
  }
}

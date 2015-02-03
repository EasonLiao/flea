package main

import (
  "fmt"
  "github.com/easonliao/flea/builtin"
  "github.com/easonliao/flea/core"
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
  fun func() error
  flag int
}

var commandsTable = map[string]cmdStruct {
  "init" : {fun : builtin.CmdInit},
  "hash-object" : {fun : builtin.CmdHashObject, flag : flagNeedSetup},
  "cat-file" : {fun : builtin.CmdCatFile, flag : flagNeedSetup},
  "status" : {fun : builtin.CmdStatus, flag : flagNeedSetup},
  "add" : {fun : builtin.CmdAdd, flag : flagNeedSetup},
  "commit" : {fun : builtin.CmdCommit, flag : flagNeedSetup},
  "branch" : {fun : builtin.CmdBranch, flag : flagNeedSetup},
  "log" : {fun : builtin.CmdLog, flag : flagNeedSetup},
  "checkout" : {fun : builtin.CmdCheckout, flag : flagNeedSetup},
}

func runBuiltin(cmd string) {
  if cmdSt, ok := commandsTable[cmd]; !ok {
    log.Fatal("Unkown command")
  } else {
    if cmdSt.flag & flagNeedSetup != 0 {
      err := core.InitFromExisting()
      if err != nil {
        fmt.Println(err.Error())
      }
    }
    err := cmdSt.fun()
    if err != nil {
      fmt.Println(err.Error())
    }
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

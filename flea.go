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
  usage func()
}

var commandsTable = map[string]cmdStruct {
  "init"        : {fun : builtin.CmdInit, usage: builtin.UsageInit},
  "hash-object" : {fun : builtin.CmdHashObject, flag : flagNeedSetup},
  "cat-file"    : {fun : builtin.CmdCatFile, flag : flagNeedSetup, usage: builtin.UsageCatFile},
  "status"      : {fun : builtin.CmdStatus, flag : flagNeedSetup, usage: builtin.UsageStatus},
  "add"         : {fun : builtin.CmdAdd, flag : flagNeedSetup, usage: builtin.UsageAdd},
  "commit"      : {fun : builtin.CmdCommit, flag : flagNeedSetup, usage: builtin.UsageCommit},
  "branch"      : {fun : builtin.CmdBranch, flag : flagNeedSetup, usage: builtin.UsageBranch},
  "log"         : {fun : builtin.CmdLog, flag : flagNeedSetup, usage: builtin.UsageLog},
  "checkout"    : {fun : builtin.CmdCheckout, flag : flagNeedSetup, usage: builtin.UsageCheckout},
  "ls-files"    : {fun : builtin.CmdLsFiles, flag : flagNeedSetup, usage: builtin.UsageLsFiles},
  "rm"          : {fun : builtin.CmdRm, flag : flagNeedSetup, usage: builtin.UsageRm},
}

func usage() {
  usage := "Here are a list of commands, see usage for specific command please use: flea <command> -h\n"
  fmt.Println(usage)
  for command, _ := range(commandsTable) {
    fmt.Printf("\t%s\n", command)
  }
  fmt.Println("")
}

func runBuiltin(cmd string) {
  if cmdSt, ok := commandsTable[cmd]; !ok {
    log.Fatal("Unkown command")
  } else {
    options := make(map[string]bool)
    for _, opt := range(os.Args[2:]) {
      options[opt] = true
    }
    if _, ok := options["-h"]; ok {
      cmdSt.usage()
    }
    if cmdSt.flag & flagNeedSetup != 0 {
      err := core.InitFromExisting()
      if err == core.ErrNoFleaDir {
        log.Fatal("Not a flea repository(or any of the parent directories):.flea")
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
    usage()
  }
}

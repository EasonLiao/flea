package core

import (
  "errors"
  "fmt"
  "io/ioutil"
  "os"
)

var (
  // ErrIO = errors.New("core: io write error")
  ErrFileNotExist = errors.New("core: file not exists")
)

func write(path string, data []byte) error {
  if err := ioutil.WriteFile(path, data, 0777); err == nil {
    return nil
  } else {
    panic(fmt.Sprintf("IO error %s while accessing %s", err, path))
  }
}

func read(path string) ([]byte, error) {
  if data, err := ioutil.ReadFile(path); err == nil {
    return data, nil
  } else if os.IsNotExist(err) {
    return nil, ErrFileNotExist
  } else {
    panic(fmt.Sprintf("IO error %s while accessing %s", err, path))
  }
}

func exists(path string) bool {
  if _, err := os.Stat(path); err == nil {
    return true
  } else if os.IsNotExist(err) {
    return false
  } else {
    panic(fmt.Sprintf("IO error %s while accessing %s", err, path))
  }
}

package core

import (
  "io/ioutil"
  "os"
  "path/filepath"
)

var testRootDir string

func init() {
  root := filepath.Join(os.TempDir(), "flea_test")
  // Removes the root directory created by last time.
  os.RemoveAll(root)
  err := os.Mkdir(root, 0777)
  if err != nil {
    panic(err.Error())
  }
  testRootDir = root
}

func createTempDir(prefix string) (name string, err error) {
  return ioutil.TempDir(testRootDir, prefix)
}

func createTempFiles(dir string, files map[string][]byte) error {
  for path, content := range(files) {
    if content == nil {
      err := os.MkdirAll(filepath.Join(dir, filepath.FromSlash(path)), 0777)
      if err != nil {
        return err
      }
    } else {
      fullpath := filepath.Join(dir, filepath.FromSlash(path))
      err := os.MkdirAll(filepath.Dir(fullpath), 0777)
      if err != nil && !os.IsExist(err) {
        return err
      }
      err = ioutil.WriteFile(fullpath, content, 0777)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func mkDir(dirname string) (string, error) {
  path := filepath.Join(testRootDir, dirname)
  return path, os.Mkdir(path, 0777)
}

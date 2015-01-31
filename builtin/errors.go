package builtin

import "errors"

var (
  ErrNotFile = errors.New("builtin: not file")
  ErrEmptyDir = errors.New("builtin: empty directory")
)

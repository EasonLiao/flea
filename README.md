# Flea
Simple Git implementation in Go. A toy project helps me to learn Go and Git.

Flea follows the design principle of Git -- A content-addressable filesystem with a VCS user interface written on top of it.

Currently Flea supports only a small subset commands of Git, you can run ```flea -h``` to inspect the list of commands it supports and the usage of these commands. The workflow is very similar to Git.

#### Initialize
```
  cd <repo>
  flea init
```

#### Commit
```
  flea add <file-path>
  flea commit -m "first commit"
```

#### Inspecting the commit history
```
  flea log
```

#### Show differences
```
  flea status
```

#### Checkout a commit snpashot from history/branch
```
  flea checkout <commit-hash>
  flea checkout master
```

### TODO
- Add branch
- More commands, e.g. revert/reset
- Show differences between files (like git diff)

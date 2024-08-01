# git-remote-foo

Playground for gittuf transport prototyping.

- provides basic `git-remote-foo` command (`main.go`), which is executed by git
  when running e.g. `git clone foo://git@git-server:git/repo`.

- *tested* with `docker compose up`, which runs:

  - A git server, which serves a bare test repo over ssh
  - A git client, which uses the transport to talk to the server (`test.sh`)
  - A log service, which outputs the transport packet trace (a log file is
    needed for debugging, because stdout/stderr are used for communication
    between transport and git parent process)

- Git pkt-line parser is vendored from `go-git` and patched to support
  gitprotocol v2 special packets (see go-git/go-git#876)

## TODO
- finalize ssh fetch (stateless-connect)
- implement ssh push
- implement curl fetch and push
- add gittuf logic
- improve testing, e.g.
  - run tests with `go test` (instead of test.sh)
  - assert for local remote git ref changes after using the transport

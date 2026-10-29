## Approach

**Code and modules.** In `read_file.go`'s `Load`, use `dt.ErrFileNotExist`. In the root go.mod drop `replace github.com/mikeschinkel/go-dt => ../go-dt` (ED-1639), then with `GOWORK=off` run `go get github.com/mikeschinkel/go-dt@v0.6.1` and `go mod tidy`. In `test/` and `examples/basic_usage`, keep only the replace pointing at go-jsonxtractr and tidy with `GOWORK=off`. A scratch copy already vets and tests clean this way.

**Main's uncommitted work.** Main has `doterr.go` modified (one line, `fmt.Fprintf(&meta, …)` in place of `meta.WriteString(fmt.Sprintf(…))`) and `.endless/config.json` staged. The modified `doterr.go` matches go-doterr's canonical copy apart from the package line, so keep it. Commit both on the task branch, then clear them from main before landing (land refuses while main holds dirty copies of files the branch changes), the same way E-2375 did: confirm the branch has them, then restore main's copies.

**justfile.** Convert the Makefile recipe for recipe (help, ensure-valid, test, test-unit, test-corpus, test-all, lint, fmt, vet, tidy, build, examples, clean, ci), keeping `GOEXPERIMENT=jsonv2`, following go-dt's justfile. Add `work`, which writes a go.work over every directory holding a go.mod (pruning `.endless`) plus `../go-dt` when it exists. go.work and go.work.sum are gitignored (add them if missing). Add a README section saying local development needs `just work`.

**Verify suite.** `.endless/tasks/e-2432/verify.toml`. Module tests use the gotest runner with an explicit command per module (`cd <dir> && GOWORK=off go test -json ./...`), since a check cannot yet name its directory (endless E-2438). go.mod facts (no root replace, go-dt v0.6.1, tidy) are generic checks.

**Tag.** After landing, give the user `v0.1.2` with its reason and a description for the current Release workflow.

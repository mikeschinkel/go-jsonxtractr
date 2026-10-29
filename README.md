# go-jsonxtractr
Go package for extracting properties from string in JSON format w/o needing to unmarshal

## Local development

The `go.mod` files require sibling packages such as
[go-dt](https://github.com/mikeschinkel/go-dt) by published version, with no
`replace` directive, so the module builds for a consumer exactly as it does
here. To build against sibling code on disk instead, you need a `go.work`
file. Create it with:

```sh
just work
```

It covers every module in this repo, plus `../go-dt` when that directory
exists. `go.work` is gitignored and is never committed; CI runs without it.

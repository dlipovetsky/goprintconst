# goprintconst

If you use Go constants as the source of truth in your build, you might end up reading their values like this:

```sh
> grep -e "MySpecial.*=.*" "foo/bar.go"
    MySpecial    = "value"
```

This operation can be error-prone: The regular expression might be incorrect, commands like `grep` or `cut` might be unavailable. By walking the Go abstract syntax tree (AST) to find constants, `goprintconst` makes the operation reliable.

## Use

```sh
> goprintconst --path foo/bar.go --name MySpecial
"value"
```

## Install

```sh
go install github.com/dlipovetsky/goprintconst
```

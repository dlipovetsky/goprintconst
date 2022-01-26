# goprintconst

If you use Go constants as the source of truth in your build, you might end up reading their values like this:

```sh
> grep -e "MySpecial.*=.*" "foo/bar.go"
    MySpecial    = "value"
```

This operation can be error-prone: The regular expression might be incorrect, commands like `grep` or `cut` might be unavailable. By walking the Go abstract syntax tree (AST) to find constants, `goprintconst` makes the operation reliable.

## Use

All constants (default):

```sh
> goprintconst --path foo/bar.go
MySpecial="value"
AnotherSpecial="foo"
ThirdSpecial="bar"
```

One constant:

```sh
> goprintconst --path foo/bar.go --names MySpecial
MySpecial="value"
```

Two constants:

```sh
> goprintconst --path foo/bar.go --names MySpecial --names AnotherSpecial
MySpecial="value"
AnotherSpecial="foo"
```

## Install

```sh
go install github.com/dlipovetsky/goprintconst
```

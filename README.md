# mockgo - mock generator for type aliased mock

mockgo generates a test mock for tye type.

## Getting started

How to install or update.

```console
$ go get -u github.com/koron/mockgo
```

Generate mock for your types.

```console
$ cd ~/go/src/github.com/your/project/package
$ mockgo -pacakge github.com/thirdparity/libarary TargetType1 TargetType2
```

### Usage

```console
$ mockgo {OPTIONS} -package {package name or relative path} [target classes...]
```

#### Options

*   `-noformat` - write mock without formatting (goimports equivalent)
*   `-output` - specify output directory (default `.`, current directory)
*   `-package` - mandatory, specify a package where types are which generate
    mock for.

#### Target classes

Where `[target classes...]` accepts two forms of name to specify type.

*   `OriginalTypename` - Mock type name will be same with `OriginalTypename`
*   `OriginalTypename:MockTypename` - Specify both original and mock type names

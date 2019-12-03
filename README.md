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

*   `-fortest` - generate mock for plain test, without `+mock` tag)
*   `-mocksuffix` - add `Mock` suffix to generated mock types
*   `-noformat` - write mock without formatting (goimports equivalent)
*   `-output` - specify output directory (default `.`, current directory)
*   `-package` - mandatory, specify a package where types are which generate
    mock for.

    starts with `./` or `../`, you can use relative path for this.

*   `-verbose` - show verbose/debug messages to stderr

#### Target classes

Where `[target classes...]` accepts two forms of name to specify type.

*   `OriginalTypename` - Mock type name will be same with `OriginalTypename`

    When `-mocksuffix` given, `OriginalTypenameMock` is used for mock type.

*   `OriginalTypename:MockTypename` - Specify both original and mock type names

    `-mocksuffix` is ignored.

## Advanced usage 

### Based component differential

for `interface` based component, `-fortest` and `-mocksuffix` will work well.

    mockgo -package ./ -outdir . -fortest -mocksuffix Interface1 Interface2

for `struct` based component, this command will work well.

    mockgo -package ../pkgA -outdir . Component1 Component2

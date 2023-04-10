# mockgo - mock generator for type aliased mock

[![GoDoc](https://godoc.org/github.com/koron/mockgo?status.svg)](https://godoc.org/github.com/koron/mockgo)
[![Actions/Go](https://github.com/koron/mockgo/workflows/Go/badge.svg)](https://github.com/koron/mockgo/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/mockgo)](https://goreportcard.com/report/github.com/koron/mockgo)

Mockgo is a tool that automatically generates mocks for unit testing based on
your types.

With the generated mock types, you can record and check method calls. Mocks can
inspect the order and arguments of the called methods, and specify the values
that methods return. These features support unit testing.

---

日本語訳/Japanese translation

mockgo はあなたの型に単体テスト用のモックを自動生成するツールです。

生成されたモック型では、メソッドの呼び出しを記録・検査できます。モックは呼び出
されるメソッドの順序や引数を検査でき、メソッドが返す値を指定できます。これらに
より単体テストを支援します。

## Getting started

How to install or update.

```console
$ go install github.com/koron/mockgo@latest
```

Generate mocks for your types: `TargetType1` and `TargetType2`

```console
$ cd ~/go/src/github.com/your/project/package
$ mockgo -pacakge github.com/thirdparity/libarary TargetType1 TargetType2
```

This generate mocks for `TargetType1` and `TargetType2` in
`"github.com/thirdparty/library"` package to current directory
(~/go/src/github.com/your/project/package).

## Usage

```console
$ mockgo {OPTIONS} -package {package name or relative path} [target classes...]
```

### Options

*   `-fortest` - generate mock for plain test, without `+mock` tag)
*   `-mocksuffix` - add `Mock` suffix to generated mock types
*   `-noformat` - write mock without formatting (goimports equivalent)
*   `-revision {num}` - mock revision 1~3. 3 is recommended, but default is 1
    for compatibility. See [mock revision](#mock-revision) for details.
*   `-noformat` - suppress goimports on generating mock code.
*   `-output {dir}` - specify output directory (default `.`, current directory)
*   `-package {path or dir}` - mandatory, packages for which there are types
    that generate mocks.

    When this starts with `./` or `../`, you can use relative path for this.

*   `-verbose` - show verbose/debug messages to stderr

### Target classes

Where `[target classes...]` accepts two forms of name to specify type.

*   `OriginalTypename` - Mock type name will be same with `OriginalTypename`

    When `-mocksuffix` given, `OriginalTypenameMock` is used for mock type.

*   `OriginalTypename:MockTypename` - Specify both original and mock type names

    `-mocksuffix` is ignored.

### Mock revision

There are three revisions of mock.

* 1 - Very simple and redundant. not recommended.

    This don't require any packages.

    * GOOD: record call parameters.
    * GOOD: specify return parameters.
    * BAD: manual check calling parameters.
    * BAD: no checks order of methods call.

* 2 - Work with expected call sequence.

    This uses a runtime `"github.com/koron/mockgo/mockrt"`.

    * GOOD: support all GOOD items in revision 1.
    * GOOD: auto check calling parameters.
    * GOOD: auto check order of methods call.
    * BAD: easily made mistakes when constructing function call sequence.

* 3 - Expected call sequnce with fault-tolerance. (recommended)

    This uses a runtime `"github.com/koron/mockgo/mockrt3"`.

    * GOOD: support all GOOD items in revision 2.
    * GOOD: fault-tolerance on constructing function call sequence.

## Type aliased mock

Usually, when using types provided by another package, you use them as they
are. However, with a type aliased mock, you first create an alias for that type
in your own package.

```go
//go:build !mock

type Foo = foo.Foo
```

This type alias will be replaced with the mock type only when building with a
specific build tag `mock`.

```go
//go:build mock

type Foo = FooMock
```

When running tests, you specify the tag to switch to the mock.

```console
$ go test -tags mock
```

This allows you to use the mock without going through an interface. This is
called a type aliased mock.

Mockgo is a command that automatically generates a mock, either FooMock or Foo
depending on whether it is a mock or not, from the source code of the
aforementioned foo.Foo type.

([Original idea from my post in Japanese](https://www.kaoriya.net/blog/2020/01/20/never-interface-only-for-tests/))

## How to check calls with mockrt3.Q

1. Create `mockrt3.Q` with `mockrt3.NewQ(*testing.T, ...)`
2. `AddCall` to add calles (`[]mockrt3.C`). A call is consist from parameter
   `P` and return values `R`. You can add calls with `NewQ` also.
3. Create a mock with `mockrt3.Q`
4. Run test target code with a mock.
5. A mock will record failures with `testing.T` when unexpected calls made by
   target code.

(TODO: Add example codes)

## Advanced usage

### Mocking `interface`

for mocking `interface` types, `-fortest` and `-mocksuffix` will work well.

```console
$ mockgo -package ./ -outdir . -revision 3 -fortest -mocksuffix Interface1 Interface2
```

This generates mock types `Interface1Mock` and `Interface2Mock` without build
tag.  `Interface1Mock` is a mock for `Interface1`, and `Interface2Mock` is for
`Interface2`.

### Mocking `struct`

for mocking `struct` types, no need special options.

```console
$ mockgo -package ../pkgA -outdir . -revision 3 Component1 Component2
```

This generates mock types `Component1` and `Component2` with `mock` build tag.
`Component1` is a mock for `pkgA.Component1`, and `Component2` is for
`pkgA.Component2`.

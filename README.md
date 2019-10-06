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
$ mockgo -outdir . -pacakge github.com/thirdparity/libarary TargetType1 TargetType2
```

### Usage

```console
$ mockgo [-noformat] -outdir {output dir} -package {package name or relative path} [target classes]
```

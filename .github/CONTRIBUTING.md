# Contributing to WISA

Contribute to our source code and to make WISA even better than it is today!

Here are the guidelines we'd like you to follow:

 - [Build Instructions](#instructions)
 - [Coding Format](#format)

## <a name="instructions"></a> Build Instructions

Building:

```go build```

Add wisa to your PATH or use it directly.

## <a name="format"></a> Coding Format

### Formatting

Format the project using `goreturns`

```goreturns .```

if you want `goreturns` write to the file for you

```goreturns -w .```

### Linting

Lint the project using `golint`

```golint ./...```

## Testing Details

#### To run tests for WISA run:

`go test -v ./...`

#### For package tests:

`go test -coverpkg=pkgname -v`

#### To visualize test coverage

`go test -coverprofile=coverage -v ./...`

then

`go tool cover -html=coverage`

This will open a browser and highlight what cases are covered and the ones that aren't

### Making Tests

In the directory containing the package you want to test, append `_test` to the filename of the file you are writing tests for. For instance, if writing tests for `package.go` the file should be named `package_test.go`
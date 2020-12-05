# WISA
Who Is Still Alive: a broken link checker

## Introduction

This is my release 0.1 project for [DPS909](https://github.com/Seneca-CDOT/topics-in-open-source-2020/wiki/release-0.1) my open source development class. It is a command-line tool for finding and reporting dead links in a file.

## Installation

### Manual
```
git clone https://github.com/1jz/wisa.git
cd wisa
go install
```

### Go
```
go get github.com/1jz/wisa
```

## Usage
type `wisa --help` for usage or:

```wisa -f [file]```

for verbose output (error logs)

```wisa -f [file] -v```

for ignoring certain URLs

```wisa -i [ignore file] -f [file]```

# Features

- Uses go routines to check for dead/broken links
- Uses `HEAD` requests for optimization
- `200-299` status codes are considered a `PASS`
- `403` emits a `WARN`
- All other codes are considered `DEAD`
- Uses flags for passing filename and for verbose output
- Coloured text using [gookit/color](https://github.com/gookit/color)
- Provide a second file to ignore certain urls


**TODO:** 

refactor the entire project as it fails to follow go concepts.

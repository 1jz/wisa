# wisa
Who Is Still Alive: a broken link checker

# Building

```go build wisa.go```

# Usage
type `wisa` for usage or: 

```wisa [file]```

# Features

- Uses go routines to check for dead/broken links
- Uses `HEAD` requests for optimization
- 200-299 status codes are considered a PASS
- 403 emits a WARN
- All other codes are considered DEAD
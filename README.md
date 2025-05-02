# Level manager for Go's slog library

`slog-leveler` manages log levels for Go's slog loggers, allowing you to instantiate and dynamically manage log levels of the default logger or component specific loggers.

## Installation

```
go get github.com/shashankram/slog-leveler
```

## Quickstart

```go
import "github.com/shashankram/slog-leveler/pkg/logger"

// Default logger
slog.Info("Hello world")
slog.Debug("This won't be printed")
logger.SetLevel(logger.DefaultComponent, slog.LevelDebug)
slog.Debug("This will be printed")
fmt.Println()

// Custom logger
fooLogger := logger.NewWithOptions("foo", logger.Options{
  Format:    logger.JSONFormat,
  AddSource: true,
})
fooLogger.Info("Hello foo")
fooLogger.Debug("This won't be printed")
logger.SetLevel("foo", slog.LevelDebug)
fooLogger.Debug("This will be printed")
fooLogger.Log(context.Background(), logger.LevelTrace, "This won't be printed")
logger.SetLevel("foo", logger.LevelTrace)
fooLogger.Log(context.Background(), logger.LevelTrace, "This will be printed")
```

## Dynamic log levels

Log levels for all loggers or a specific component logger can be dynamically changed by either using `logger.SetLevel` or by using the `logger.HTTPLevelHandler` HTTP handler function.

Log level can be set to `error|warn|info|debug|trace`.

### HTTP handler

The following examples demonstrate using the HTTP log level handler to dynamically change and view log levels.
It assumes that the HTTP server is running on `localhost:15000`.

Update the log level of all loggers:
```
curl -X POST 'localhost:15000/logging?level=<level>'
```

Update the log level of a single component logger:
```
curl -X POST 'localhost:15000/logging?foo=debug'
```

Update the log level of a multiple component loggers:
```
curl -X POST 'localhost:15000/logging?foo=debug&bar=trace'
```

Display the current log level of all loggers:
```console
$ curl -X POST 'localhost:15000/logging'
current log levels:
---
bar: trace
default: info
foo: info
```

## Examples

- [Quickstart](/examples/quickstart/quickstart.go)
- [HTTP handler](/examples/dynamic/logging.go)
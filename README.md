`astilog` goal is to provide an implementation of [astikit.CompleteLogger](https://github.com/asticode/go-astikit/blob/master/logger.go#L8) interface while restricting dependencies to [astikit](https://github.com/asticode/go-astikit) only.

It doesn't provide a global logger anymore since it should be a project-based decision and libs need an instance of `astikit.CompleteLogger` provided to them anyway.

It doesn't use `logrus` anymore since most of its features were not used. As a result `astilog` is now much simpler and quicker since it only implements a subset of `logrus` features.

# Installation

Run the following command:

```
$ go get github.com/asticode/go-astilog
```

# Usage

## Create a logger using explicit configuration

```go
// Create logger
l := astilog.New(astilog.Configuration{
    AppName: "myapp",
    Format: astilog.FormatText,
    Level: astilog.LevelWarn,
    Out: astilog.OutStderr,
    Source: true,
})

// Make sure to close the logger properly
defer l.Close()
```

## Create a logger using flags only

```go
// Parse flags
flag.Parse()

// Create logger
l := astilog.NewFromFlags()

// Make sure to close the logger properly
defer l.Close()
```

## Log stuff

```go
l.Debug("this is a log message")
l.Infof("this is a %s message", "log")
```

## Log stuff using the context

```go
// Add field to the context
ctx = astilog.ContextWithField(ctx, "key", "value")

// Log with context
l.DebugC(ctx, "this is a log message")
l.InfoCf(ctx, "this is a %s message", "log")
```

## Outputs
### Log to file

Set the `Filename` option to your file path.

### Log to syslog

Set the `Out` option to `syslog` or `astilog.OutSyslog` if you're setting it in GO.

### Log to stderr

Set the `Out` option to `stderr` or `astilog.OutStderr` if you're setting it in GO.

### Log to stdout

Set the `Out` option to `stdout` or `astilog.OutStdout` if you're setting it in GO.
package astilog

import "flag"

// Flags
var (
	AppName         = flag.String("logger-app-name", "", "the logger app name")
	Filename        = flag.String("logger-filename", "", "the logger filename")
	Format          = flag.String("logger-format", "", "the logger format")
	Level           = flag.String("logger-level", "", "the logger level")
	MessageKey      = flag.String("logger-message-key", "", "the logger message key")
	Out             = flag.String("logger-out", "", "the logger out")
	Source          = flag.Bool("logger-source", false, "if true, then source is added to fields")
	TimestampFormat = flag.String("logger-timestamp-format", "", "the logger timestamp format")
	Verbose         = flag.Bool("v", false, "if true, then log level is debug")
)

// Formats
const (
	FormatJSON = "json"
	FormatText = "text"
)

// Levels
const (
	LevelDebug = "debug"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelInfo  = "info"
	LevelWarn  = "warn"
)

// Outs
const (
	OutStderr = "stderr"
	OutStdout = "stdout"
	OutSyslog = "syslog"
)

// Configuration represents the configuration of the logger
type Configuration struct {
	AppName         string `toml:"app_name"`
	Filename        string `toml:"filename"`
	Format          string `toml:"format"`
	Level           string `toml:"level"`
	MessageKey      string `toml:"message_key"`
	Out             string `toml:"out"`
	Source          bool   `toml:"source"`
	TimestampFormat string `toml:"timestamp_format"`
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() (c Configuration) {
	c = Configuration{
		AppName:         *AppName,
		Filename:        *Filename,
		Format:          *Format,
		Level:           *Level,
		MessageKey:      *MessageKey,
		Out:             *Out,
		Source:          *Source,
		TimestampFormat: *TimestampFormat,
	}
	if *Verbose {
		c.Level = LevelDebug
	}
	return
}

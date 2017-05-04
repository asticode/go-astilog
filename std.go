package astilog

// Global logger
var gb = NopLogger()

// SetLogger sets the global logger
func SetLogger(l Logger) {
	gb = l
}

// Global logger shortcuts
func Debug(v ...interface{})                 { gb.Debug(v...) }
func Debugf(format string, v ...interface{}) { gb.Debugf(format, v...) }
func Info(v ...interface{})                  { gb.Info(v...) }
func Infof(format string, v ...interface{})  { gb.Infof(format, v...) }
func Warn(v ...interface{})                  { gb.Warn(v...) }
func Warnf(format string, v ...interface{})  { gb.Warnf(format, v...) }
func Error(v ...interface{})                 { gb.Error(v...) }
func Errorf(format string, v ...interface{}) { gb.Errorf(format, v...) }
func Fatal(v ...interface{})                 { gb.Fatal(v...) }
func Fatalf(format string, v ...interface{}) { gb.Fatalf(format, v...) }

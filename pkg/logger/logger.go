package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// lock is a global mutex lock to gain control of logrus.<SetLevel|SetOutput>
var lock = sync.Mutex{}

var Formatter = new(logrus.TextFormatter)

// log executes a logging function with the specified output stream and mutex protection
func log(output *os.File, fn func()) {
	lock.Lock()
	logrus.SetOutput(output)
	fn()
	lock.Unlock()
}

// SetLevel sets the standard logger level to the specified logrus.Level
func SetLevel(level logrus.Level) {
	lock.Lock()
	logrus.SetLevel(level)
	lock.Unlock()
}

// SetLevelDebug sets the standard logger level to Debug
func SetLevelDebug() {
	SetLevel(logrus.DebugLevel)
}

// SetLevelInfo sets the standard logger level to Info
func SetLevelInfo() {
	SetLevel(logrus.InfoLevel)
}

// SetLevelWarn sets the standard logger level to Warn
func SetLevelWarn() {
	SetLevel(logrus.WarnLevel)
}

// SetLevelError sets the standard logger level to Error
func SetLevelError() {
	SetLevel(logrus.ErrorLevel)
}

// Trace logs a message at level Trace to stdout.
func Trace(args ...interface{}) {
	log(os.Stdout, func() { logrus.Trace(args...) })
}

// Tracef logs a message at level Trace to stdout.
func Tracef(format string, args ...interface{}) {
	log(os.Stdout, func() { logrus.Tracef(format, args...) })
}

// Traceln logs a message at level Trace to stdout.
func Traceln(args ...interface{}) {
	log(os.Stdout, func() { logrus.Traceln(args...) })
}

// Debug logs a message at level Debug to stdout.
func Debug(args ...interface{}) {
	log(os.Stdout, func() { logrus.Debug(args...) })
}

// Debugf logs a message at level Debug to stdout.
func Debugf(format string, args ...interface{}) {
	log(os.Stdout, func() { logrus.Debugf(format, args...) })
}

// Debugln logs a message at level Debug to stdout.
func Debugln(args ...interface{}) {
	log(os.Stdout, func() { logrus.Debugln(args...) })
}

// Info logs a message at level Info to stdout.
func Info(args ...interface{}) {
	log(os.Stdout, func() { logrus.Info(args...) })
}

// Infof logs a message at level Info to stdout.
func Infof(format string, args ...interface{}) {
	log(os.Stdout, func() { logrus.Infof(format, args...) })
}

// Infoln logs a message at level Info to stdout.
func Infoln(args ...interface{}) {
	log(os.Stdout, func() { logrus.Infoln(args...) })
}

// Warn logs a message at level Warn to stdout.
func Warn(args ...interface{}) {
	log(os.Stdout, func() { logrus.Warn(args...) })
}

// Warnf logs a message at level Warn to stdout.
func Warnf(format string, args ...interface{}) {
	log(os.Stdout, func() { logrus.Warnf(format, args...) })
}

// Warnln logs a message at level Warn to stdout.
func Warnln(args ...interface{}) {
	log(os.Stdout, func() { logrus.Warnln(args...) })
}

// Error logs a message at level Error to stderr.
func Error(args ...interface{}) {
	log(os.Stderr, func() { logrus.Error(args...) })
}

// Errorf logs a message at level Error to stderr.
func Errorf(format string, args ...interface{}) {
	log(os.Stderr, func() { logrus.Errorf(format, args...) })
}

// Errorln logs a message at level Error to stderr.
func Errorln(args ...interface{}) {
	log(os.Stderr, func() { logrus.Errorln(args...) })
}

// Fatal logs a message at level Fatal to stderr then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	log(os.Stderr, func() { logrus.Fatal(args...) })
}

// Fatalf logs a message at level Fatal to stderr.
func Fatalf(format string, args ...interface{}) {
	log(os.Stderr, func() { logrus.Fatalf(format, args...) })
}

// Fatalln logs a message at level Fatal to stderr.
func Fatalln(args ...interface{}) {
	log(os.Stderr, func() { logrus.Fatalln(args...) })
}

// Panic logs a message at level Panic to stderr; calls panic() after logging.
func Panic(args ...interface{}) {
	log(os.Stderr, func() { logrus.Panic(args...) })
}

// Panicf logs a message at level Panic to stderr.
func Panicf(format string, args ...interface{}) {
	log(os.Stderr, func() { logrus.Panicf(format, args...) })
}

// Panicln logs a message at level Panic to stderr.
func Panicln(args ...interface{}) {
	log(os.Stderr, func() { logrus.Panicln(args...) })
}

func init() {
	// Setup logger defaults
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	Formatter.DisableTimestamp = true
	logrus.SetFormatter(Formatter)
	logrus.SetOutput(os.Stdout) // Set output to stdout; set to stderr by default
	logrus.SetLevel(logrus.DebugLevel)
}

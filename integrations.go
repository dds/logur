package logur

import (
	"fmt"
	"strings"

	"github.com/goph/logur/internal/keyvals"
)

/*
gRPC uses a globally configured logger that implements the google.golang.org/grpc/grpclog.LoggerV2 interface.

GRPCV2Logger is an adapter around logur-compatible loggers,
so you can easily wire the logging library of your choice into gRPC:

	package main

	import (
		"github.com/goph/logur"
		"google.golang.org/grpc/grpclog"
	)

	func main() {
		logger := logur.NewNoopLogger() // choose an actual implementation
		grpclog.SetLoggerV2(logur.NewGRPCV2Logger(logger))
	}
*/
type GRPCV2Logger struct {
	logger       Logger
	levelEnabler LevelEnabler
}

// NewGRPCV2Logger returns a new V2 gRPC logger.
func NewGRPCV2Logger(logger Logger) *GRPCV2Logger {
	l := &GRPCV2Logger{
		logger: logger,
	}

	if levelEnabler, ok := logger.(LevelEnabler); ok {
		l.levelEnabler = levelEnabler
	}

	return l
}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (l *GRPCV2Logger) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *GRPCV2Logger) Infoln(args ...interface{}) {
	l.logger.Info(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (l *GRPCV2Logger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (l *GRPCV2Logger) Warning(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (l *GRPCV2Logger) Warningln(args ...interface{}) {
	l.logger.Warn(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (l *GRPCV2Logger) Warningf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (l *GRPCV2Logger) Error(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (l *GRPCV2Logger) Errorln(args ...interface{}) {
	l.logger.Error(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (l *GRPCV2Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (l *GRPCV2Logger) Fatal(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (l *GRPCV2Logger) Fatalln(args ...interface{}) {
	l.logger.Error(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (l *GRPCV2Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

// V reports whether verbosity level l is at least the requested verbose level.
func (l *GRPCV2Logger) V(level int) bool {
	if l.levelEnabler == nil {
		return true
	}

	if level == 3 { // fatal level
		level = 2
	}

	// grpc log doesn't have trace and debug levels
	return l.levelEnabler.LevelEnabled(Level(level + 2))
}

/*
KitLogger turns a logur logger into a Go-Kit logger.

	package main

	import (
		"github.com/goph/logur"
	)

	func main() {
		logger := logur.NewNoopLogger() // choose an actual implementation
		kitlogger := logur.NewKitLogger(logger)

		// inject the logger somewhere
	}
*/
type KitLogger struct {
	logFuncs       map[string]LogFunc
	defaultLogFunc LogFunc
}

// NewKitLogger returns a new go-kit logger.
func NewKitLogger(logger Logger) *KitLogger {
	l := &KitLogger{
		logFuncs: map[string]LogFunc{
			"trace":   logger.Trace,
			"debug":   logger.Debug,
			"info":    logger.Info,
			"warn":    logger.Warn,
			"warning": logger.Warn,
			"error":   logger.Error,
		},
		defaultLogFunc: logger.Info,
	}

	return l
}

func (l *KitLogger) Log(kvs ...interface{}) error {
	if len(kvs)%2 == 1 {
		kvs = append(kvs, "(MISSING)")
	}

	fields := keyvals.ToMap(kvs)

	logFunc := l.defaultLogFunc

	if lf, ok := l.logFuncs[strings.ToLower(fmt.Sprintf("%s", fields["level"]))]; ok {
		delete(fields, "level")

		logFunc = lf
	}

	var msg string
	if m, ok := fields["msg"]; ok {
		delete(fields, "msg")
		msg = fmt.Sprintf("%s", m)
	}

	logFunc(msg, fields)

	return nil
}
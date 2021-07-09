// Package contract provides interfaces for a common logging backend
package contract

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type fieldImpl struct {
	key   string
	value interface{}
}

func (f *fieldImpl) Key() string        { return f.key }
func (f *fieldImpl) Value() interface{} { return f.value }

// NewField provides a simple shortcut function to create a struct that
// satisfies the Field interface
func NewField(key string, value interface{}) Field {
	return &fieldImpl{
		key:   key,
		value: value,
	}
}

type fieldCollection []Field

func (f fieldCollection) String() string {
	var s []string
	for _, item := range f {
		s = append(s, fmt.Sprintf("%s=%q", item.Key(), fmt.Sprintf("%v", item.Value())))
	}
	return strings.Join(s, ", ")
}

type logger struct {
	outWriter          io.Writer
	errWriter          io.Writer
	minLevel           Level
	isInDevEnvironment bool
	disableAllWrites   bool
	name               string
	fields             fieldCollection
}

func newLogger(outWriter, errWriter io.Writer, minLevel Level) *logger {
	return &logger{
		outWriter:          outWriter,
		errWriter:          errWriter,
		minLevel:           minLevel,
		isInDevEnvironment: true,
	}
}

func (l *logger) Named(name string) Logger {
	if len(l.name) > 0 {
		name = l.name + "." + name
	}
	name = strings.ToUpper(name)

	return &logger{
		outWriter:          l.outWriter,
		errWriter:          l.errWriter,
		minLevel:           l.minLevel,
		disableAllWrites:   l.disableAllWrites,
		isInDevEnvironment: l.isInDevEnvironment,
		name:               name,
		fields:             l.fields,
	}
}

func (l *logger) With(fields ...Field) Logger {
	return &logger{
		outWriter:          l.outWriter,
		errWriter:          l.errWriter,
		minLevel:           l.minLevel,
		disableAllWrites:   l.disableAllWrites,
		isInDevEnvironment: l.isInDevEnvironment,
		name:               l.name,
		fields:             append(l.fields, fields...),
	}
}

func (l *logger) Sync() error {
	return nil
}

func (l *logger) log(level Level, msg string, fields fieldCollection) {
	if l.disableAllWrites {
		return
	}

	if level < l.minLevel {
		return
	}

	timeStr := time.Now().UTC().Format(time.RFC3339Nano)
	for len(timeStr) < 35 {
		timeStr += " "
	}

	var levelStr string
	switch level {
	case DebugLevel:
		levelStr = "DEBUG "
	case InfoLevel:
		levelStr = "INFO  "
	case WarnLevel:
		levelStr = "WARN  "
	case ErrorLevel:
		levelStr = "ERROR "
	case DPanicLevel:
		levelStr = "DPANIC"
	case PanicLevel:
		levelStr = "PANIC "
	case FatalLevel:
		levelStr = "FATAL "
	}

	var nameStr string
	if len(l.name) > 0 {
		nameStr = l.name + ": "
	}

	s := fmt.Sprintf("%s %s %s%s (%s)\n", timeStr, levelStr, nameStr, msg, append(l.fields, fields...).String())
	if _, err := fmt.Fprintf(l.outWriter, s); err != nil {
		_, _ = fmt.Fprintf(l.errWriter, "logger: writing of message %q failed due to: %v", s, err)
	}
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields)
}

func (l *logger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields)
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields)
}

func (l *logger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields)
}

func (l *logger) DPanic(msg string, fields ...Field) {
	l.log(DPanicLevel, msg, fields)
	if l.isInDevEnvironment {
		panic(msg)
	}
}

func (l *logger) Panic(msg string, fields ...Field) {
	l.log(PanicLevel, msg, fields)
	panic(msg)
}

func (l *logger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields)
	os.Exit(1)
}

// StdImplOption is an option function for NewStdImpl and MustNewStdImpl
type StdImplOption func(*logger) error

// DisableLogWrites fully disables the logger instance. No message is written
// to the OutWriter and ErrWriter.
func DisableLogWrites() StdImplOption {
	return func(log *logger) error {
		log.disableAllWrites = true
		return nil
	}
}

// OutWriter lets you set the destination for the default logging output
func OutWriter(w io.Writer) StdImplOption {
	return func(log *logger) error {
		log.outWriter = w
		return nil
	}
}

// ErrWriter lets you set the destination for errors that happened inside the
// standard logger itself (which at the moment can only happen if the OutWriter
// returns an error on writing)
func ErrWriter(w io.Writer) StdImplOption {
	return func(log *logger) error {
		log.errWriter = w
		return nil
	}
}

// MinLevel sets the minimum needed level for messages that get written to
// OutWriter. Messages below this level will get discarded
func MinLevel(minLevel Level) StdImplOption {
	return func(log *logger) error {
		log.minLevel = minLevel
		return nil
	}
}

// IsInDevEnvironment influences whether or not DPanicLevel panics or not
func IsInDevEnvironment(isInDevEnvironment bool) StdImplOption {
	return func(log *logger) error {
		log.isInDevEnvironment = isInDevEnvironment
		return nil
	}
}

// NewStdImpl creates a new Logger using only Go's standard library. This can
// be useful for packages that want to include liblog, and provide logging
// output by default, but not overstuff their library with third-party
// dependencies.
func NewStdImpl(opts ...StdImplOption) (Logger, error) {
	log := newLogger(os.Stdout, os.Stderr, DebugLevel)
	for _, opt := range opts {
		if err := opt(log); err != nil {
			return nil, err
		}
	}
	return log, nil
}

// MustNewStdImpl is like NewStdImpl, but panics if StdImplOption cannot be
// applied.
func MustNewStdImpl(opts ...StdImplOption) Logger {
	l, err := NewStdImpl(opts...)
	if err != nil {
		panic(err)
	}
	return l
}

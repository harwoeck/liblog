package contract

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type fieldCollection []Field

func (f fieldCollection) String() string {
	var s []string
	for _, item := range f {
		s = append(s, fmt.Sprintf("%s=%q", item.Key(), fmt.Sprintf("%v", item.Value())))
	}
	r := strings.Join(s, ", ")
	if len(r) > 0 {
		return "(" + r + ")"
	}
	return ""
}

type logger struct {
	outWriter        io.Writer
	errWriter        io.Writer
	wd               string
	optMinLevel      Level
	optInDev         bool
	optDisableWrites bool
	name             string
	fields           fieldCollection
}

func newLogger(outWriter, errWriter io.Writer, minLevel Level) *logger {
	wd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(errWriter, "liblog/contract.std: failed to get working directory\n")
	}

	return &logger{
		outWriter:   outWriter,
		errWriter:   errWriter,
		wd:          wd,
		optMinLevel: minLevel,
		optInDev:    true,
	}
}

func (l *logger) clone() *logger {
	return &logger{
		outWriter:        l.outWriter,
		errWriter:        l.errWriter,
		wd:               l.wd,
		optMinLevel:      l.optMinLevel,
		optInDev:         l.optInDev,
		optDisableWrites: l.optDisableWrites,
		name:             l.name,
		fields:           l.fields,
	}
}

func (l *logger) Named(name string) Logger {
	clone := l.clone()
	if len(clone.name) == 0 {
		clone.name = name
	} else {
		clone.name += "." + name
	}
	return clone
}

func (l *logger) With(fields ...Field) Logger {
	clone := l.clone()
	clone.fields = append(clone.fields, fields...)
	return clone
}

func (l *logger) Sync() error {
	return nil
}

func (l *logger) log(level Level, msg string, fields fieldCollection) {
	if l.optDisableWrites {
		return
	}

	if level < l.optMinLevel {
		return
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05.999 Z07:00")
	for len(timeStr) < 30 {
		timeStr += " "
	}

	levelStr := level.String()
	for len(levelStr) < 6 {
		levelStr += " "
	}

	var nameStr string
	if len(l.name) > 0 {
		nameStr = l.name + " "
	}

	var callerStr string
	_, frameF, frameL, defined := runtime.Caller(2)
	if !defined {
		_, _ = fmt.Fprintf(l.errWriter, "liblog/contract.std: failed to get caller\n")
	} else {
		if len(l.wd) > 0 && l.wd != "/" {
			frameF = strings.TrimPrefix(frameF, l.wd)
			frameF = strings.TrimPrefix(frameF, "/")
		}

		callerStr = fmt.Sprintf("%s:%d ", frameF, frameL)
	}

	s := fmt.Sprintf("%s %s %s%s%s %s\n", timeStr, levelStr, nameStr, callerStr, msg, fields.String())

	if _, err := fmt.Fprintf(l.outWriter, s); err != nil {
		_, _ = fmt.Fprintf(l.errWriter, "liblog/contract.std: writing of message %q failed due to: %v", s, err)
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

// FormatToError is a public helper function that converts a msg and fields pair
// into an combined error. Intended for Logger.ErrorReturn
func FormatToError(name string, callerSkip int, msg string, fields ...Field) error {
	var nameStr string
	if len(name) > 0 {
		nameStr = name + ": "
	}

	wd, _ := os.Getwd()

	var caller string
	_, frameF, frameL, defined := runtime.Caller(callerSkip + 1)
	if defined {
		if len(wd) > 0 && wd != "/" {
			frameF = strings.TrimPrefix(frameF, wd)
			frameF = strings.TrimPrefix(frameF, "/")
		}

		caller = fmt.Sprintf("%s:%d: ", frameF, frameL)
	}

	fieldsStr := fieldCollection(fields).String()
	if len(fieldsStr) > 0 {
		fieldsStr = " " + fieldsStr
	}

	return fmt.Errorf("%s%s%s%s", nameStr, caller, msg, fieldsStr)
}

func (l *logger) ErrorReturn(msg string, fields ...Field) error {
	l.log(ErrorLevel, msg, fields)
	return FormatToError(l.name, 1, msg, fields...)
}

func (l *logger) DPanic(msg string, fields ...Field) {
	l.log(DPanicLevel, msg, fields)
	if l.optInDev {
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

// StdOption is an option function for NewStd and MustNewStd
type StdOption func(*logger) error

// DisableLogWrites fully disables the logger instance. No message is written
// to the OutWriter and ErrWriter.
func DisableLogWrites() StdOption {
	return func(log *logger) error {
		log.optDisableWrites = true
		return nil
	}
}

// OutWriter lets you set the destination for the default logging output
func OutWriter(w io.Writer) StdOption {
	return func(log *logger) error {
		log.outWriter = w
		return nil
	}
}

// ErrWriter lets you set the destination for errors that happened inside the
// standard logger itself (which at the moment can only happen if the OutWriter
// returns an error on writing)
func ErrWriter(w io.Writer) StdOption {
	return func(log *logger) error {
		log.errWriter = w
		return nil
	}
}

// MinLevel sets the minimum needed level for messages that get written to
// OutWriter. Messages below this level will get discarded
func MinLevel(minLevel Level) StdOption {
	return func(log *logger) error {
		log.optMinLevel = minLevel
		return nil
	}
}

// IsInDevEnvironment influences whether or not DPanicLevel panics or not
func IsInDevEnvironment(isInDevEnvironment bool) StdOption {
	return func(log *logger) error {
		log.optInDev = isInDevEnvironment
		return nil
	}
}

// NewStd creates a new Logger using only Go's standard library. This can
// be useful for packages that want to include liblog, and provide logging
// output by default, but not overstuff their library with third-party
// dependencies.
func NewStd(opts ...StdOption) (Logger, error) {
	log := newLogger(os.Stdout, os.Stderr, DebugLevel)
	for _, opt := range opts {
		if err := opt(log); err != nil {
			return nil, err
		}
	}
	return log, nil
}

// MustNewStd is like NewStd, but panics if one of the StdOption cannot be
// applied.
func MustNewStd(opts ...StdOption) Logger {
	l, err := NewStd(opts...)
	if err != nil {
		panic(err)
	}
	return l
}

// Package contract provides interfaces for a common logging backend
package contract

// A Field is a marshaling operation used to add a key-value pair to a logger's
// context
type Field interface {
	// Key returns the fields key part
	Key() string
	// Value returns the fields value part
	Value() interface{}
}

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

// Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

// String returns a readable representation of Level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case DPanicLevel:
		return "DPANIC"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides leveled, structured logging. It is an abstract interface for
// different logging backends.
type Logger interface {
	// Named adds a new path segment to the logger's name. Segments are joined by
	// periods. By default, Loggers are unnamed.
	Named(name string) Logger
	// With creates a child logger and adds structured context to it. Fields added
	// to the child don't affect the parent, and vice versa.
	With(fields ...Field) Logger
	// Sync flushes any buffered log entries and should be called by applications
	// before exiting
	Sync() error
	// Debug logs a message at DebugLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Debug(msg string, fields ...Field)
	// Info logs a message at InfoLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Info(msg string, fields ...Field)
	// Warn logs a message at WarnLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Warn(msg string, fields ...Field)
	// Error logs a message at ErrorLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Error(msg string, fields ...Field)
	// ErrorReturn logs a message at ErrorLevel exactly like the Error function, but
	// additionally returns an error object, containing the provided information.
	ErrorReturn(msg string, fields ...Field) error
	// DPanic logs a message at DPanicLevel. The message includes any fields
	// passed at the log site, as well as any fields accumulated on the logger.
	//
	// If the logger is in development mode, it then panics (DPanic means
	// "development panic"). This is useful for catching errors that are
	// recoverable, but shouldn't ever happen.
	DPanic(msg string, fields ...Field)
	// Panic logs a message at PanicLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	//
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic(msg string, fields ...Field)
	// Fatal logs a message at FatalLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	//
	// The logger then calls os.Exit(1), even if logging at FatalLevel is
	// disabled.
	Fatal(msg string, fields ...Field)
}

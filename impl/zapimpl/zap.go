package zapimpl

import (
	"reflect"
	"unsafe"

	"github.com/harwoeck/liblog/contract"
	"go.uber.org/zap"
)

func NewZapImpl(log *zap.Logger) contract.Logger {
	return newZapImpl(log.WithOptions(zap.AddCallerSkip(1)))
}

func castFields(fields []contract.Field) []zap.Field {
	zf := make([]zap.Field, 0)
	for _, item := range fields {
		zf = append(zf, zap.Any(item.Key(), item.Value()))
	}
	return zf
}

func getZapLoggerName(log *zap.Logger) string {
	rs := reflect.ValueOf(log).Elem()
	rf := rs.FieldByName("name")
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	return rf.String()
}

type impl struct {
	log    *zap.Logger
	name   string
	fields []contract.Field
}

func newZapImpl(log *zap.Logger) *impl {
	return &impl{
		log:  log,
		name: getZapLoggerName(log),
	}
}

func (i *impl) Named(name string) contract.Logger {
	l := newZapImpl(i.log.Named(name))
	l.fields = i.fields
	return l
}

func (i *impl) With(fields ...contract.Field) contract.Logger {
	l := newZapImpl(i.log.With(castFields(fields)...))
	l.fields = append(i.fields, fields...)
	return l
}

func (i *impl) Sync() error {
	return i.log.Sync()
}

func (i *impl) Debug(msg string, fields ...contract.Field) {
	i.log.Debug(msg, castFields(fields)...)
}

func (i *impl) Info(msg string, fields ...contract.Field) {
	i.log.Info(msg, castFields(fields)...)
}

func (i *impl) Warn(msg string, fields ...contract.Field) {
	i.log.Warn(msg, castFields(fields)...)
}

func (i *impl) Error(msg string, fields ...contract.Field) {
	i.log.Error(msg, castFields(fields)...)
}

func (i *impl) ErrorReturn(msg string, fields ...contract.Field) error {
	i.log.Error(msg, castFields(fields)...)
	return contract.FormatToError(i.name, 1, msg, append(i.fields, fields...)...)
}

func (i *impl) DPanic(msg string, fields ...contract.Field) {
	i.log.DPanic(msg, castFields(fields)...)
}

func (i *impl) Panic(msg string, fields ...contract.Field) {
	i.log.Panic(msg, castFields(fields)...)
}

func (i *impl) Fatal(msg string, fields ...contract.Field) {
	i.log.Fatal(msg, castFields(fields)...)
}

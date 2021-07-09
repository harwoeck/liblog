package zapimpl

import (
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

type impl struct {
	log *zap.Logger
}

func newZapImpl(log *zap.Logger) contract.Logger {
	return &impl{
		log: log,
	}
}

func (i *impl) Named(name string) contract.Logger {
	return newZapImpl(i.log.Named(name))
}

func (i *impl) With(fields ...contract.Field) contract.Logger {
	return newZapImpl(i.log.With(castFields(fields)...))
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

func (i *impl) DPanic(msg string, fields ...contract.Field) {
	i.log.DPanic(msg, castFields(fields)...)
}

func (i *impl) Panic(msg string, fields ...contract.Field) {
	i.log.Panic(msg, castFields(fields)...)
}

func (i *impl) Fatal(msg string, fields ...contract.Field) {
	i.log.Fatal(msg, castFields(fields)...)
}

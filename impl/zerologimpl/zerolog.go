package zerologimpl

import (
	"github.com/harwoeck/liblog/contract"
	"github.com/rs/zerolog"
)

func NewZerologImpl(log *zerolog.Logger, inDevMode bool) contract.Logger {
	return newImpl(log, inDevMode)
}

func castFields(fields []contract.Field) interface{} {
	m := make(map[string]interface{})
	for _, item := range fields {
		m[item.Key()] = item.Value()
	}
	return m
}

type impl struct {
	log       *zerolog.Logger
	inDevMode bool
	name      string
	fields    []contract.Field
}

func newImpl(log *zerolog.Logger, inDevMode bool) *impl {
	return &impl{
		log:       log,
		inDevMode: inDevMode,
		name:      "",
		fields:    nil,
	}
}

func (i *impl) Named(name string) contract.Logger {
	if len(i.name) > 0 {
		name = i.name + "." + name
	}

	sub := i.log.With().Str("name", name).Logger()
	l := newImpl(&sub, i.inDevMode)
	l.name = name
	l.fields = i.fields

	return l
}

func (i *impl) With(fields ...contract.Field) contract.Logger {
	sub := i.log.With().Fields(castFields(fields)).Logger()
	l := newImpl(&sub, i.inDevMode)
	l.name = i.name
	l.fields = append(l.fields, fields...)
	return l
}

func (i *impl) Sync() error {
	return nil
}

func (i *impl) Debug(msg string, fields ...contract.Field) {
	i.log.Debug().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) Info(msg string, fields ...contract.Field) {
	i.log.Info().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) Warn(msg string, fields ...contract.Field) {
	i.log.Warn().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) Error(msg string, fields ...contract.Field) {
	i.log.Error().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) ErrorReturn(msg string, fields ...contract.Field) error {
	i.Error(msg, fields...)
	return contract.FormatToError(i.name, 1, msg, append(i.fields, fields...)...)
}

func (i *impl) DPanic(msg string, fields ...contract.Field) {
	if i.inDevMode {
		i.Panic(msg, fields...)
	} else {
		i.Error(msg, fields...)
	}
}

func (i *impl) Panic(msg string, fields ...contract.Field) {
	i.log.Panic().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) Fatal(msg string, fields ...contract.Field) {
	i.log.Fatal().Fields(castFields(fields)).Msg(msg)
}

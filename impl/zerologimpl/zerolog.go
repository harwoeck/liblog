package zerologimpl

import (
	"github.com/harwoeck/liblog"
	"github.com/rs/zerolog"
)

func NewZerologImpl(log *zerolog.Logger, inDevMode bool) liblog.Logger {
	return newImpl(log, inDevMode)
}

func castFields(fields []liblog.Field) interface{} {
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
	fields    []liblog.Field
}

func newImpl(log *zerolog.Logger, inDevMode bool) *impl {
	return &impl{
		log:       log,
		inDevMode: inDevMode,
		name:      "",
		fields:    nil,
	}
}

func (i *impl) Named(name string) liblog.Logger {
	if len(i.name) > 0 {
		name = i.name + "." + name
	}

	sub := i.log.With().Str("name", name).Logger()
	l := newImpl(&sub, i.inDevMode)
	l.name = name
	l.fields = i.fields

	return l
}

func (i *impl) With(fields ...liblog.Field) liblog.Logger {
	sub := i.log.With().Fields(castFields(fields)).Logger()
	l := newImpl(&sub, i.inDevMode)
	l.name = i.name
	l.fields = append(l.fields, fields...)
	return l
}

func (i *impl) Sync() error {
	return nil
}

func (i *impl) Debug(msg string, fields ...liblog.Field) {
	i.log.Debug().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) Info(msg string, fields ...liblog.Field) {
	i.log.Info().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) Warn(msg string, fields ...liblog.Field) {
	i.log.Warn().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) Error(msg string, fields ...liblog.Field) {
	i.log.Error().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) ErrorReturn(msg string, fields ...liblog.Field) error {
	i.Error(msg, fields...)
	return liblog.FormatToError(i.name, 1, msg, append(i.fields, fields...)...)
}

func (i *impl) DPanic(msg string, fields ...liblog.Field) {
	if i.inDevMode {
		i.Panic(msg, fields...)
	} else {
		i.Error(msg, fields...)
	}
}

func (i *impl) Panic(msg string, fields ...liblog.Field) {
	i.log.Panic().Fields(castFields(fields)).Msg(msg)
}

func (i *impl) Fatal(msg string, fields ...liblog.Field) {
	i.log.Fatal().Fields(castFields(fields)).Msg(msg)
}

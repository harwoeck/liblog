package zerologimpl

import (
	"os"
	"testing"
	"time"

	"github.com/harwoeck/liblog"
	"github.com/rs/zerolog"
)

var field = liblog.NewField

func TestNewZerologImpl(t *testing.T) {
	z := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	i := NewZerologImpl(&z, false)

	z.Info().Str("foo", "bar").Msg("Hello World")
	i.Info("Hello World", field("foo", "bar"))
	i.Named("service").With(field("foo", "bar")).Info("Hello World")
	i.Named("service").With(field("foo", "bar")).Named("db").Info("upserted")
}

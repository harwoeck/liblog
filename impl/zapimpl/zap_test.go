package zapimpl

import (
	"testing"

	"github.com/harwoeck/liblog"
	"go.uber.org/zap"
)

var field = liblog.NewField

func TestNewZapImpl(t *testing.T) {
	z, _ := zap.NewDevelopment()
	i := NewZapImpl(z)

	z.Info("Hello World", zap.String("foo", "bar"))
	i.Info("Hello World", field("foo", "bar"))
	i.Named("service").With(field("foo", "bar")).Info("Hello World")
	i.Named("service").With(field("foo", "bar")).Named("db").Info("upserted")
}

package zapimpl

import (
	"github.com/harwoeck/liblog/contract"
	"go.uber.org/zap"
	"testing"
)

func TestNewZapImpl(t *testing.T) {
	z, _ := zap.NewDevelopment()
	i := NewZapImpl(z)

	z.Info("Hello World", zap.String("foo", "bar"))
	i.Info("Hello World", contract.NewField("foo", "bar"))
	i.Named("service").With(contract.NewField("foo", "bar")).Info("Hello World")
	i.Named("service").With(contract.NewField("foo", "bar")).Named("db").Info("upserted")
}

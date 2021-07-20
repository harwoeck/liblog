# liblog

*liblog* is an interface/contract for logging backends. It should be used by public libraries and packages that want to give their user's control over structured and leveled logging output.

[![Go Reference](https://pkg.go.dev/badge/github.com/harwoeck/liblog/contract.svg)](https://pkg.go.dev/github.com/harwoeck/liblog/contract)

### Advantages

- 🟢 Users can provide their own logging stack __**and**__ get detailed package-level logging
- 🟢 *liblog* provides implementations for well-known logging backends (see `/impl`)

### Getting started

- `go get github.com/harwoeck/liblog/contract` - provides the *liblog* interface, and an implementation that only relies on Go's standard library
- Use `contract.NewStd(...StdOption)` by default inside your package, but give users the option to specify their own logging implementation (`contract.Logger`)

### Usage

- Extensively use all provided options, like `Named()`, `With()` and appropriate logging levels, to provide the best log results/experience.
- Use `field()` (or `contract.NewField()`) to generate structured elements for your logs (see Tips #1)

### Tips

1. Create a file `logfield.go` and add this `fieldImpl`. You can now use `field()` instead of the longer `contract.NewField()` inside your package to create structured logging elements
    ```go
    type fieldImpl struct {
        key   string
        value interface{}
    }
    
    func (f *fieldImpl) Key() string        { return f.key }
    func (f *fieldImpl) Value() interface{} { return f.value }
    
    func field(key string, value interface{}) contract.Field {
        return &fieldImpl{
            key:   key,
            value: value,
        }
    }
    ```

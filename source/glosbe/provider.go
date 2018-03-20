// Copyright © 2018 Trevor N. Suarez (Rican7)

package glosbe

import (
	"net/http"

	flag "github.com/ogier/pflag"

	"github.com/Rican7/define/registry"
	"github.com/Rican7/define/source"
)

type config struct{}

type provider struct{}

func init() {
	registry.Register(Name, registry.RegisterFunc(register))
}

func register(*flag.FlagSet) (registry.SourceProvider, registry.Configuration) {
	return &provider{}, nil
}

func (p *provider) Provide(config registry.Configuration) (source.Source, error) {
	return New(http.Client{}), nil
}

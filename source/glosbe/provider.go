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

// JSONKey defines the JSON key used for the provider
const JSONKey = "GlosbeAPI"

func init() {
	registry.Register(registry.RegisterFunc(register))
}

func register(*flag.FlagSet) (registry.SourceProvider, registry.Configuration) {
	return &provider{}, &config{}
}

func (c *config) JSONKey() string {
	return JSONKey
}

func (p *provider) Name() string {
	return Name
}

func (p *provider) Provide(conf registry.Configuration) (source.Source, error) {
	return New(http.Client{}), nil
}

// Copyright © 2023 Trevor N. Suarez (Rican7)

package freedictionaryapi

import (
	"net/http"

	flag "github.com/ogier/pflag"

	"github.com/Rican7/define/registry"
	"github.com/Rican7/define/source"
)

// RequiredConfigError represents an error when a required configuration key is
// missing or invalid.
type RequiredConfigError struct {
	Key string
}

type config struct{}

type provider struct{}

// JSONKey defines the JSON key used for the provider
const JSONKey = "FreeDictionaryAPI"

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

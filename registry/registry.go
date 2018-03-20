// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package registry provides a registry for sources and their providers to
// integrate into the source list
package registry

import (
	"sync"

	flag "github.com/ogier/pflag"

	"github.com/Rican7/define/source"
)

// SourceProvider defines the interface for providers of sources.
type SourceProvider interface {
	// Provide returns a source based on a given configuration.
	Provide(Configuration) (source.Source, error)
}

// Configuration defines a generic SourceProvider's configuration structure.
//
// Implementations may wish to implement the json.Marshaler and
// json.Unmarshaler interfaces to customize their JSON representations.
type Configuration interface {
	// JSONKey returns the JSON key that should be used when marshalling and
	// unmarshalling the configuration into a global/shared configuration JSON
	// representation. Due to the nature of the global/shared status, this key
	// should be "unique", so as not to overwrite/clash with other provider
	// configurations.
	JSONKey() string
}

// RegisterFunc is the function that allows SourceProviders to define and
// expose their configuration structure to the registry, so that sources can be
// provided with a dynamically initialized configuration.
type RegisterFunc func(*flag.FlagSet) (SourceProvider, Configuration)

var (
	configured    sync.Once
	registrations = make(map[string]RegisterFunc)

	providers = make(map[string]SourceProvider)
)

// Register makes a source provider available by the provided name.
func Register(name string, registerFunc RegisterFunc) {
	registrations[name] = registerFunc
}

// ConfigureProviders configures the providers and returns a map of their names
// as keys and their configurations as values.
//
// This is intended to be called ONLY by the registry owner.
// TODO: Prevent external calls somehow?
func ConfigureProviders(flags *flag.FlagSet) map[string]Configuration {
	confs := make(map[string]Configuration)

	configured.Do(func() {
		for name, registerFunc := range registrations {
			providers[name], confs[name] = registerFunc(flags)
		}
	})

	return confs
}

// Provide TODO
func Provide(name string, conf Configuration) (source.Source, error) {
	return providers[name].Provide(conf)
}

// Providers returns a list of the names of the registered providers.
func Providers() []string {
	var list []string

	for name := range providers {
		list = append(list, name)
	}

	return list
}

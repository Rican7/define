// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package registry provides a registry for sources and their providers to
// integrate into the source list
package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	flag "github.com/ogier/pflag"

	"github.com/Rican7/define/source"
)

// SourceProvider defines the interface for providers of sources.
type SourceProvider interface {
	// Name returns a printable user-friendly name to refer to the source by.
	Name() string

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

// DynamicConfiguration defines a generic SourceProvider's configuration
// structure that allows for a dynamic loading mechanism.
type DynamicConfiguration interface {
	Configuration
	json.Unmarshaler

	// Finalize is a method called on a configuration when loading is completed.
	//
	// The intent is to be able signal to the configuration that its been loaded
	// so that it can perform any necessary post-load processes, such as
	// validation, normalization, or having values fall-back to defaults.
	Finalize()
}

// RegisterFunc is the function that allows SourceProviders to define and
// expose their configuration structure to the registry, so that sources can be
// provided with a dynamically initialized configuration.
type RegisterFunc func(*flag.FlagSet) (SourceProvider, Configuration)

var (
	configured    sync.Once
	finalized     sync.Once
	registrations = make([]RegisterFunc, 0)

	providers = make(map[Configuration]SourceProvider)
)

// Register makes a source provider available by the provided name.
func Register(registerFunc RegisterFunc) {
	registrations = append(registrations, registerFunc)
}

// ConfigureProviders configures the providers and returns a map of their names
// as keys and their configurations as values.
//
// This is intended to be called ONLY by the registry owner.
// TODO: Prevent external calls somehow?
func ConfigureProviders(flags *flag.FlagSet) map[string]Configuration {
	confs := make(map[string]Configuration)

	configured.Do(func() {
		for _, registerFunc := range registrations {
			provider, conf := registerFunc(flags)

			if nil == provider || nil == conf {
				panic("register func returned nil values")
			}

			providers[conf], confs[conf.JSONKey()] = provider, conf
		}
	})

	return confs
}

// Finalize takes a number of configurations and marks them as loaded, if they
// support a DynamicConfiguration signaling.
//
// This is intended to be called ONLY by the registry owner.
// TODO: Prevent external calls somehow?
func Finalize(confs ...Configuration) {
	finalized.Do(func() {
		for _, conf := range confs {
			if dynamicConf, ok := conf.(DynamicConfiguration); ok {
				dynamicConf.Finalize()
			}
		}
	})
}

// Provide takes a configuration and calls the associated source providers
// Provide function to provide a source.
func Provide(conf Configuration) (source.Source, error) {
	provider := providers[conf]

	src, err := provider.Provide(conf)

	if nil != err {
		err = fmt.Errorf("source %q failed to initialize with error: %s", provider.Name(), err)
	}

	return src, err
}

// ProvidePreferred takes a preferred provider key (that aligns with the value
// returned by the Configuration.JSONKey method) and a list of configurations,
// and provides the matching source if possible, but will fall back to another
// source if the preferred source returns an error when trying to provide it.
func ProvidePreferred(preferredProvider string, confs []Configuration) (source.Source, error) {
	var src source.Source
	var err error

	if len(confs) < 1 {
		return nil, errors.New("no configurations available to provide a source")
	}

	for _, providerConf := range confs {
		if src == nil || nil != err || preferredProvider == providerConf.JSONKey() {
			iSrc, iErr := Provide(providerConf)

			if nil != iSrc && nil == iErr {
				src, err = iSrc, iErr
			}
		}
	}

	return src, err
}

// Providers returns a map of the source configurations as keys and their
// corresponding providers as values.
func Providers() map[Configuration]SourceProvider {
	provs := make(map[Configuration]SourceProvider)

	for conf, provider := range providers {
		provs[conf] = provider
	}

	return provs
}

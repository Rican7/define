// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package config provides types and handling for configuration values used in
// the application.
package config

import (
	"os"
	"strconv"

	flag "github.com/ogier/pflag"

	"github.com/imdario/mergo"
)

// Configuration defines the application's configuration structure
type Configuration struct {
	IndentationSize uint
	PreferredSource string
}

// initializeCommandLineConfig initializes the command line configuration from
// a list of arguments that are parsed as flags.
func initializeCommandLineConfig(flags *flag.FlagSet, args []string) (Configuration, error) {
	var conf Configuration

	// Define our flags
	flags.UintVar(&conf.IndentationSize, "indent-size", 0, "The number of spaces to indent output by")
	flags.StringVar(&conf.PreferredSource, "preferred-source", "", "The preferred source to use, if available")

	err := flags.Parse(args)

	return conf, err
}

// initializeEnvironmentConfig initializes the environment configuration from
// the application's environment.
func initializeEnvironmentConfig() Configuration {
	var conf Configuration

	if val, err := strconv.ParseUint(os.Getenv("DEFINE_APP_INDENT_SIZE"), 10, 0); nil == err {
		conf.IndentationSize = uint(val)
	}

	conf.PreferredSource = os.Getenv("DEFINE_APP_PREFERRED_SOURCE")

	return conf
}

// mergeConfigurations merges multiple configurations values together, from left
// to right argument position, by filling any of the left arguments zero-values
// with any non-zero-values from the right.
func mergeConfigurations(confs ...Configuration) (Configuration, error) {
	var merged Configuration

	for _, conf := range confs {
		if err := mergo.Merge(&merged, conf); nil != err {
			return merged, err
		}
	}

	return merged, nil
}

// NewFromRuntime builds a Configuration by merging values from multiple
// different sources. It accepts a Configuration containing default values to
// fill in any empty/blank configuration values found when merging from the
// different sources.
//
// The merging of values from different sources will take this priority:
// 1. Command line arguments
// 2. Environment variables
// 3. Passed in default values
func NewFromRuntime(flags *flag.FlagSet, defaults Configuration) (Configuration, error) {
	var conf Configuration
	var err error

	commandLineConfig, err := initializeCommandLineConfig(flags, os.Args[1:])

	if nil == err {
		conf, err = mergeConfigurations(
			commandLineConfig,
			initializeEnvironmentConfig(),
			defaults,
		)
	}

	return conf, err
}

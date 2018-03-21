// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package config provides types and handling for configuration values used in
// the application.
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"
	flag "github.com/ogier/pflag"

	"github.com/imdario/mergo"
)

// Configuration defines the application's configuration structure
type Configuration struct {
	ConfigFileLocation string

	IndentationSize uint
	PreferredSource string
}

// initializeCommandLineConfig initializes the command line configuration from
// a list of arguments that are parsed as flags.
func initializeCommandLineConfig(flags *flag.FlagSet, args []string) (Configuration, error) {
	var conf Configuration

	// Define our flags
	flags.StringVarP(&conf.ConfigFileLocation, "config-file", "c", "", "The location of the config file to use")
	flags.UintVar(&conf.IndentationSize, "indent-size", 0, "The number of spaces to indent output by")
	flags.StringVar(&conf.PreferredSource, "preferred-source", "", "The preferred source to use, if available")

	err := flags.Parse(args)

	return conf, err
}

func initializeFileConfig(fileLocation string) (Configuration, error) {
	var conf Configuration

	// If we can expand the location, do so
	if expanded, err := homedir.Expand(fileLocation); nil == err {
		fileLocation = expanded
	}

	fileContents, err := ioutil.ReadFile(fileLocation)

	if nil != err {
		return conf, err
	}

	err = json.Unmarshal(fileContents, &conf)

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
// 2. A loaded config file, if available
// 3. Environment variables
// 4. Passed in default values
func NewFromRuntime(flags *flag.FlagSet, defaults Configuration) (Configuration, error) {
	var conf Configuration
	var err error

	var commandLineConfig Configuration
	var fileConfig Configuration

	commandLineConfig, err = initializeCommandLineConfig(flags, os.Args[1:])

	if nil == err {
		configFileLocation := commandLineConfig.ConfigFileLocation

		if "" == configFileLocation && "" != defaults.ConfigFileLocation {
			// If we can expand the location, do so
			if expanded, err := homedir.Expand(defaults.ConfigFileLocation); nil == err {
				defaults.ConfigFileLocation = expanded
			}

			// If we haven't passed a config file flag, and our default exists
			if _, err := os.Stat(defaults.ConfigFileLocation); !os.IsNotExist(err) {
				// Set our location to the default, since it exists
				// (if there are problems reading the file, we'll handle later)
				configFileLocation = defaults.ConfigFileLocation
			}
		}

		// If we have a config file to load
		if "" != configFileLocation {
			fileConfig, err = initializeFileConfig(configFileLocation)

			if nil != err {
				err = errors.New("Error reading config file")
			}
		}
	}

	if nil == err {
		conf, err = mergeConfigurations(
			commandLineConfig,
			fileConfig,
			initializeEnvironmentConfig(),
			defaults,
		)
	}

	return conf, err
}

// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package action provides types for determining the action intended to be
// performed by the application.
package action

import (
	flag "github.com/ogier/pflag"
)

// List of actions to perform.
const (
	DefineWord Type = iota
	PrintConfig
	DebugConfig
	ListSources
	PrintVersion
)

// Type defines the type of action intended for the app to perform.
type Type uint

// Action defines an intended action for the app to perform.
type Action struct {
	flagSet *flag.FlagSet
	flag    struct {
		printConfig  bool
		debugConfig  bool
		listSources  bool
		printVersion bool
	}
}

// Setup sets up a lazy-valued action based on a given flag set.
//
// NOTE: The passed flag set will have to be parsed before the action can be
// used. This parsing is intended to not be performed in this
func Setup(flags *flag.FlagSet) *Action {
	var act Action

	// Define our flags
	flags.BoolVar(&act.flag.printConfig, "print-config", false, "To print the current configuration")
	flags.BoolVar(&act.flag.debugConfig, "debug-config", false, "To print debug info about the configuration")
	flags.BoolVar(&act.flag.listSources, "list-sources", false, "To print the available sources")
	flags.BoolVar(&act.flag.printVersion, "version", false, "To print the app's version info")

	// Pass our flagset, so we can be diligent about parse checking later
	act.flagSet = flags

	return &act
}

// validateState validates that the action state is valid and panics if not.
func (a *Action) validateState() {
	if !a.flagSet.Parsed() {
		panic("action is unusable; flags haven't been parsed")
	}
}

// Type returns the action type to perform.
func (a *Action) Type() Type {
	a.validateState()

	switch {
	case a.flag.printConfig:
		return PrintConfig
	case a.flag.debugConfig:
		return DebugConfig
	case a.flag.listSources:
		return ListSources
	case a.flag.printVersion:
		return PrintVersion
	default:
		return DefineWord
	}
}

// Copyright © 2018 Trevor N. Suarez (Rican7)

// A command-line dictionary (thesaurus) app, written in Go.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Rican7/define/internal/action"
	"github.com/Rican7/define/internal/config"
	defineio "github.com/Rican7/define/internal/io"
	"github.com/Rican7/define/internal/io/printer"
	"github.com/Rican7/define/internal/version"
	"github.com/Rican7/define/registry"
	"github.com/Rican7/define/source"
	flag "github.com/ogier/pflag"

	_ "github.com/Rican7/define/source/freedictionaryapi"
	"github.com/Rican7/define/source/oxford"
	_ "github.com/Rican7/define/source/webster"
)

const (
	// Configuration defaults
	defaultConfigFileLocation = "~/.define.conf.json"
	defaultIndentationSize    = 2
	defaultPreferredSource    = oxford.JSONKey
)

var (
	stdErrWriter = defineio.NewPanicWriter(os.Stderr, defaultIndentationSize)
	stdOutWriter = defineio.NewPanicWriter(os.Stdout, defaultIndentationSize)

	flags *flag.FlagSet
	act   *action.Action
	conf  config.Configuration
	src   source.Source
)

func init() {
	var err error

	flags = flag.NewFlagSet(version.AppName, flag.ContinueOnError)
	flags.SetOutput(stdErrWriter)
	flags.Usage = func() {
		printUsage(stdErrWriter)
		quit(2)
	}

	act = action.Setup(flags)

	// Configure our registered providers
	providerConfs := registry.ConfigureProviders(flags)
	var providerConfsList []registry.Configuration

	if len(providerConfs) < 1 {
		handleError(fmt.Errorf("no registered source providers"))
	}

	for _, providerConf := range providerConfs {
		providerConfsList = append(providerConfsList, providerConf)
	}

	conf, err = config.NewFromRuntime(flags, providerConfs, defaultConfigFileLocation, config.Configuration{
		IndentationSize: defaultIndentationSize,
		PreferredSource: defaultPreferredSource,
	})

	// Re-initialize our writers once we have our indentation size configuration
	stdErrWriter = defineio.NewPanicWriter(os.Stderr, conf.IndentationSize)
	stdOutWriter = defineio.NewPanicWriter(os.Stdout, conf.IndentationSize)
	flags.SetOutput(stdErrWriter)

	// Finalize our configurations
	registry.Finalize(providerConfsList...)

	handleError(err)

	if conf.Source != "" {
		if providerConf, exists := providerConfs[conf.Source]; exists {
			src, err = registry.Provide(providerConf)
		} else {
			handleError(fmt.Errorf("provider/source %q does not exist", conf.Source))
		}
	} else {
		src, err = registry.ProvidePreferred(conf.PreferredSource, providerConfsList)
	}

	// Make sure our flags are parsed before entering main
	handleError(err, flags.Parse(os.Args[1:]))
}

func handleSourceError(source string, err ...error) {
	for _, e := range err {
		if e == nil {
			continue
		}

		msg := e.Error()

		if len(msg) > 1 {
			// Format the message
			msg = strings.ToTitle(msg[:1]) + msg[1:]

			stdErrWriter.IndentWrites(func(writer *defineio.PanicWriter) {
				if source != "" {
					sourceMessage := fmt.Sprintf("Source %q encountered an error.", source)

					writer.WriteNewLine()
					writer.WriteStringLine(sourceMessage)
				}

				writer.WritePaddedStringLine(msg, 1)
			})
		}

		quit(1)
	}
}

func handleError(err ...error) {
	handleSourceError("", err...)
}

func quit(code int) {
	os.Exit(code)
}

func printConfig() {
	encoded, err := json.MarshalIndent(conf, "", "    ")

	handleError(err)

	stdOutWriter.WriteStringLine(string(encoded))
}

func printSources() {
	var sourceStrings []string

	for conf, source := range registry.Providers() {
		sourceStrings = append(sourceStrings, fmt.Sprintf("%q (%s)", source.Name(), conf.JSONKey()))
	}

	sort.Strings(sourceStrings)

	stdOutWriter.IndentWrites(func(writer *defineio.PanicWriter) {
		writer.WritePaddedStringLine("Available sources:", 1)

		for i, source := range sourceStrings {
			writer.WriteStringLine(fmt.Sprintf("%d. %s", i+1, source))
		}

		writer.WriteNewLine()
	})
}

func printVersion() {
	stdOutWriter.WriteStringLine(version.Printable())
}

func printUsage(writer *defineio.PanicWriter) {
	writer.IndentWrites(func(w *defineio.PanicWriter) {
		flags.SetOutput(w)

		w.WritePaddedStringLine(fmt.Sprintf("Usage: %s [<options>...] <word>", version.AppName), 1)

		w.WriteStringLine("Options:")
		flags.PrintDefaults()
		w.WriteNewLine()
	})
}

func defineWord(word string) {
	results, err := src.Define(word)

	handleSourceError(src.Name(), err, source.ValidateDictionaryResults(word, results))

	resultPrinter := printer.NewResultPrinter(stdOutWriter)

	resultPrinter.PrintDictionaryResults(results)
	resultPrinter.PrintSourceName(src)
}

func main() {
	// Get the word from our first non-flag argument
	word := flags.Arg(0)

	// Decide what to perform
	switch act.Type() {
	case action.PrintConfig:
		printConfig()
	case action.ListSources:
		printSources()
	case action.PrintVersion:
		printVersion()
	case action.DefineWord:
		fallthrough
	default:
		if word == "" {
			// Show our usage
			printUsage(stdOutWriter)
			quit(1)
		} else {
			defineWord(word)
		}
	}
}

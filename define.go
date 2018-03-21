// Copyright © 2018 Trevor N. Suarez (Rican7)

// A command-line dictionary (thesaurus) app, written in Go.
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/Rican7/define/internal/action"
	"github.com/Rican7/define/internal/config"
	defineio "github.com/Rican7/define/internal/io"
	"github.com/Rican7/define/registry"
	"github.com/Rican7/define/source"
	flag "github.com/ogier/pflag"

	_ "github.com/Rican7/define/source/glosbe"
	"github.com/Rican7/define/source/oxford"
	_ "github.com/Rican7/define/source/webster"
)

const (
	appName = "define"

	// Configuration defaults
	defaultConfigFileLocation = "~/.define.conf"
	defaultIndentationSize    = 2
	defaultPreferredSource    = oxford.JSONKey
)

var (
	stdErrWriter = defineio.NewPanicWriter(os.Stderr)
	stdOutWriter = defineio.NewPanicWriter(os.Stdout)

	flags *flag.FlagSet
	act   *action.Action
	conf  config.Configuration
	src   source.Source
)

func init() {
	var err error

	flags = flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.Usage = func() {
		printUsage(stdErrWriter, defaultIndentationSize)
		quit(2)
	}
	flags.SetOutput(stdErrWriter.Writer())

	act = action.Setup(flags)

	// Configure our registered providers
	providerConfs := registry.ConfigureProviders(flags)

	if len(providerConfs) < 1 {
		handleError(fmt.Errorf("No registered source providers"))
	}

	conf, err = config.NewFromRuntime(flags, providerConfs, defaultConfigFileLocation, config.Configuration{
		IndentationSize: defaultIndentationSize,
		PreferredSource: defaultPreferredSource,
	})

	handleError(err)

	var preferredProviderConfig registry.Configuration

	if "" != conf.PreferredSource {
		if providerConf, ok := providerConfs[conf.PreferredSource]; ok {
			preferredProviderConfig = providerConf
		} else {
			handleError(fmt.Errorf("Preferred provider/source %q does not exist", conf.PreferredSource))
		}
	} else {
		for _, providerConf := range providerConfs {
			preferredProviderConfig = providerConf
			break
		}
	}

	src, err = registry.Provide(preferredProviderConfig)

	handleError(err)

	// Make sure our flags are parsed before entering main
	handleError(flags.Parse(os.Args[1:]))
}

func main() {
	// Get the word from our first non-flag argument
	word := flags.Arg(0)

	// Decide what to perform
	switch act.Type() {
	case action.PrintConfig:
		printConfig()
	case action.DefineWord:
		fallthrough
	default:
		if "" == word {
			// Show our usage
			printUsage(stdOutWriter, conf.IndentationSize)
			quit(1)
		} else {
			defineWord(word)
		}
	}
}

func handleError(err ...error) {
	for _, e := range err {
		if nil != e {
			stdErrWriter.IndentWrites(conf.IndentationSize, func(writer *defineio.PanicWriter) {
				writer.WriteNewLine()
				writer.WriteStringLine(e.Error())
				writer.WriteNewLine()
			})

			quit(1)
		}
	}
}

func quit(code int) {
	os.Exit(code)
}

func printConfig() {
	encoded, err := json.MarshalIndent(conf, "", "    ")

	handleError(err)

	stdOutWriter.WriteStringLine(string(encoded))
}

func printUsage(writer *defineio.PanicWriter, indentSize uint) {
	writer.IndentWrites(indentSize, func(w *defineio.PanicWriter) {
		flags.SetOutput(w.Writer())

		writer.WriteNewLine()
		writer.WriteStringLine(fmt.Sprintf("Usage: %s [<options>...] <word>", appName))
		writer.WriteNewLine()

		writer.WriteStringLine("Options:")
		flags.PrintDefaults()
		writer.WriteNewLine()
	})
}

func defineWord(word string) {
	result, err := src.Define(word)

	handleError(err, source.ValidateResult(result))

	printResult(result)
	printSourceName(src)
}

func printSourceName(src source.Source) {
	stdOutWriter.IndentWrites(conf.IndentationSize, func(writer *defineio.PanicWriter) {
		text := fmt.Sprintf("Results provided by: %q", src.Name())
		separatorSize := int(math.Min(float64(60), float64(len(text))))

		writer.WriteNewLine()
		writer.WriteStringLine(strings.Repeat("-", separatorSize))
		writer.WriteStringLine(text)
		writer.WriteNewLine()
	})
}

func printResult(result source.Result) {
	stdOutWriter.IndentWrites(conf.IndentationSize, func(writer *defineio.PanicWriter) {
		writer.WriteNewLine()
		writer.WriteStringLine(getHeader(result))
		writer.WriteNewLine()

		for _, entry := range result.Entries() {
			if entryHeader := getEntryHeader(result, entry); "" != entryHeader {
				writer.WriteNewLine()
				writer.WriteNewLine()
				writer.WriteStringLine(entryHeader)
			}

			writer.IndentWrites(conf.IndentationSize, func(writer *defineio.PanicWriter) {

				if wordEntry, isWordEntry := entry.(source.WordEntry); isWordEntry && "" != wordEntry.Category() {
					writer.WriteNewLine()
					writer.WriteStringLine(fmt.Sprintf("(%s)", wordEntry.Category()))
					writer.WriteNewLine()
				}

				for senseIndex, sense := range entry.Senses() {
					prefix := fmt.Sprintf("%d. ", senseIndex+1)

					for defIndex, definition := range sense.Definitions() {
						// Change the prefix after the first definition
						if 0 < defIndex {
							prefix = " - "
						}

						writer.WriteStringLine(prefix + definition)
					}

					writer.IndentWrites(uint(len(prefix)), func(writer *defineio.PanicWriter) {
						for _, examples := range sense.Examples() {
							writer.WriteStringLine(fmt.Sprintf("%q", examples))
						}

						for _, notes := range sense.Notes() {
							writer.WriteStringLine(fmt.Sprintf("[%s]", notes))
						}
					})

					writer.IndentWrites(conf.IndentationSize, func(writer *defineio.PanicWriter) {
						for _, subSense := range sense.Subsenses() {
							prefix := " - "

							for _, definition := range subSense.Definitions() {
								writer.WriteStringLine(prefix + definition)
							}

							writer.IndentWrites(uint(len(prefix)), func(writer *defineio.PanicWriter) {
								if len(subSense.Examples()) > 0 {
									writer.WriteStringLine(fmt.Sprintf("%q", subSense.Examples()[0]))
								}
							})
						}
					})
				}

				if etymologyEntry, ok := entry.(source.EtymologyEntry); ok {
					if 0 < len(etymologyEntry.Etymologies()) {
						writer.WriteNewLine()
						writer.WriteStringLine("Origin")
						writer.WriteNewLine()

						for _, etymology := range etymologyEntry.Etymologies() {
							writer.WriteStringLine(etymology)
						}

						writer.WriteNewLine()
					}
				}

				if thesaurusEntry, ok := entry.(source.ThesaurusEntry); ok {
					if 0 < len(thesaurusEntry.Synonyms()) {
						writer.WriteNewLine()
						writer.WriteStringLine("Synonyms")
						writer.WriteNewLine()

						writer.WriteStringLine(strings.Join(thesaurusEntry.Synonyms(), " ; "))

						writer.WriteNewLine()
					}

					if 0 < len(thesaurusEntry.Antonyms()) {
						writer.WriteNewLine()
						writer.WriteStringLine("Antonyms")
						writer.WriteNewLine()

						writer.WriteStringLine(strings.Join(thesaurusEntry.Antonyms(), " ; "))

						writer.WriteNewLine()
					}
				}
			})
		}

		writer.WriteNewLine()
	})
}

func getHeader(result source.Result) string {
	header := result.Headword()

	firstEntry := result.Entries()[0]

	if "" != firstEntry.Pronunciation() || (isSameWord(result, firstEntry) && "" != firstEntry.Pronunciation()) {
		header = fmt.Sprintf("%s  /%s/", header, firstEntry.Pronunciation())
	}

	return header
}

func getEntryHeader(result source.Result, entry source.DictionaryEntry) string {
	var header string

	if wordEntry, isWordEntry := entry.(source.WordEntry); isWordEntry && !isSameWord(result, entry) {
		if "" != entry.Pronunciation() {
			header = fmt.Sprintf("%s  /%s/", wordEntry.Word(), entry.Pronunciation())
		} else {
			header = wordEntry.Word()
		}
	}

	return header
}

func isSameWord(result source.Result, entry source.DictionaryEntry) bool {
	wordEntry, isWordEntry := entry.(source.WordEntry)

	return isWordEntry && wordEntry.Word() == result.Headword()
}

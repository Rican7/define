// A command-line dictionary (thesaurus) app, written in Go.
//
// Copyright © 2018 Trevor N. Suarez (Rican7)
package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"

	defineio "github.com/Rican7/define/io"
	"github.com/Rican7/define/source"
	"github.com/Rican7/define/source/glosbe"
)

var indentSize = 2

func main() {
	var word string

	if len(os.Args) >= 2 {
		word = os.Args[1]
	}

	src := glosbe.New(http.Client{})

	result, err := src.Define(word)

	handleError(err, source.ValidateResult(result))

	stdOutWriter := defineio.NewPanicWriter(os.Stdout)

	printResult(result, stdOutWriter)
	printSourceName(src, stdOutWriter)
}

func handleError(err ...error) {
	for _, e := range err {
		if nil != e {
			errorAndQuit(e)
		}
	}
}

func errorAndQuit(err error) {
	defineio.NewPanicWriter(os.Stderr).IndentWrites(indentSize, func(writer *defineio.PanicWriter) {
		writer.WriteNewLine()
		writer.WriteStringLine(err.Error())
		writer.WriteNewLine()
	})

	os.Exit(1) // TODO: Error codes?
}

func printSourceName(src source.Source, writer *defineio.PanicWriter) {
	writer.IndentWrites(indentSize, func(writer *defineio.PanicWriter) {
		text := fmt.Sprintf("Results provided by: %q", src.Name())
		separatorSize := int(math.Min(float64(60), float64(len(text))))

		writer.WriteNewLine()
		writer.WriteStringLine(strings.Repeat("-", separatorSize))
		writer.WriteStringLine(text)
		writer.WriteNewLine()
	})
}

func printResult(result source.Result, writer *defineio.PanicWriter) {
	writer.IndentWrites(indentSize, func(writer *defineio.PanicWriter) {
		writer.WriteNewLine()
		writer.WriteStringLine(getHeader(result))
		writer.WriteNewLine()

		for _, entry := range result.Entries() {
			if entryHeader := getEntryHeader(result, entry); "" != entryHeader {
				writer.WriteNewLine()
				writer.WriteNewLine()
				writer.WriteStringLine(entryHeader)
			}

			writer.IndentWrites(indentSize, func(writer *defineio.PanicWriter) {

				if wordEntry, isWordEntry := entry.(source.WordEntry); isWordEntry && "" != wordEntry.Category() {
					writer.WriteNewLine()
					writer.WriteStringLine(fmt.Sprintf("(%s)", wordEntry.Category()))
					writer.WriteNewLine()
				}

				for senseIndex, sense := range entry.Senses() {
					prefix := fmt.Sprintf("%d. ", senseIndex+1)

					for _, definition := range sense.Definitions() {
						writer.WriteStringLine(prefix + definition)
					}

					writer.IndentWrites(len(prefix), func(writer *defineio.PanicWriter) {
						for _, examples := range sense.Examples() {
							writer.WriteStringLine(fmt.Sprintf("%q", examples))
						}

						for _, notes := range sense.Notes() {
							writer.WriteStringLine(fmt.Sprintf("[%s]", notes))
						}
					})

					writer.IndentWrites(indentSize, func(writer *defineio.PanicWriter) {
						for _, subSense := range sense.Subsenses() {
							prefix := " - "

							for _, definition := range subSense.Definitions() {
								writer.WriteStringLine(prefix + definition)
							}

							writer.IndentWrites(len(prefix), func(writer *defineio.PanicWriter) {
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

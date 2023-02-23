// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package printer provides types and methods to encapsulate consistent printing
// of structures.
package printer

import (
	"fmt"
	"math"
	"strings"

	defineio "github.com/Rican7/define/internal/io"
	"github.com/Rican7/define/source"
)

const (
	etymologyHeader = "Origin"
	synonymHeader   = "Synonyms"
	antonymHeader   = "Antonyms"
)

// ResultPrinter is a printer for source.Result structures.
type ResultPrinter struct {
	out *defineio.PanicWriter
}

// NewResultPrinter creates a new ResultPrinter.
func NewResultPrinter(out *defineio.PanicWriter) *ResultPrinter {
	return &ResultPrinter{out: out}
}

// PrintSourceName prints the name of a source.Source.
func (p *ResultPrinter) PrintSourceName(src source.Source) {
	p.out.IndentWrites(func(writer *defineio.PanicWriter) {
		text := fmt.Sprintf("Results provided by: %q", src.Name())
		separatorSize := int(math.Min(float64(60), float64(len(text))))

		writer.WriteNewLine()
		writer.WriteStringLine(strings.Repeat("-", separatorSize))
		writer.WriteStringLine(text)
		writer.WriteNewLine()
	})
}

// PrintDictionaryResults prints a list of dictionary results
func (p *ResultPrinter) PrintDictionaryResults(results []source.DictionaryResult) {
	p.out.IndentWrites(func(writer *defineio.PanicWriter) {
		var lastWord string

		for _, result := range results {
			resultHeader := getHeader(result)
			writer.WritePaddedStringLine(resultHeader, 1)

			var lastEntryHeader string
			for _, entry := range result.Entries {
				if entryHeader := getEntryHeader(resultHeader, lastEntryHeader, lastWord, entry); entryHeader != "" {
					writer.WriteNewLine()
					writer.WriteNewLine()
					writer.WriteStringLine(entryHeader)

					lastEntryHeader = entryHeader
				}

				writer.IndentWrites(func(writer *defineio.PanicWriter) {
					printDictionaryEntry(writer, entry)
				})

				lastWord = entry.Word
			}

			writer.WriteNewLine()
		}
	})
}

// PrintSearchResults prints a list of search results
func (p *ResultPrinter) PrintSearchResults(results []string) {
	p.out.IndentWrites(func(writer *defineio.PanicWriter) {
		for index, result := range results {
			writer.WriteStringLine(fmt.Sprintf("%d. %s", index+1, result))
		}
	})
}

func printDictionaryEntry(writer *defineio.PanicWriter, entry source.DictionaryEntry) {
	if entry.LexicalCategory != "" {
		writer.WritePaddedStringLine(fmt.Sprintf("(%s)", entry.LexicalCategory), 1)
	}

	for senseIndex, sense := range entry.Senses {
		prefix := fmt.Sprintf("%d. ", senseIndex+1)

		for defIndex, definition := range sense.Definitions {
			// Change the prefix after the first definition
			if 0 < defIndex {
				prefix = " - "
			}

			writer.WriteStringLine(prefix + definition)
		}

		writer.IndentWritesBy(uint(len(prefix)), func(writer *defineio.PanicWriter) {
			for _, examples := range sense.Examples {
				writer.WriteStringLine(examples.String())
			}

			for _, notes := range sense.Notes {
				writer.WriteStringLine(fmt.Sprintf("[%s]", notes))
			}
		})

		writer.IndentWrites(func(writer *defineio.PanicWriter) {
			for _, subSense := range sense.SubSenses {
				prefix := " - "

				for _, definition := range subSense.Definitions {
					writer.WriteStringLine(prefix + definition)
				}

				writer.IndentWritesBy(uint(len(prefix)), func(writer *defineio.PanicWriter) {
					if len(subSense.Examples) > 0 {
						writer.WriteStringLine(subSense.Examples[0].String())
					}
				})
			}
		})
	}

	printEtymologies(writer, entry)
	printThesaurusValues(writer, entry.ThesaurusValues)
}

func printEtymologies(writer *defineio.PanicWriter, entry source.DictionaryEntry) {
	if 0 < len(entry.Etymologies) {
		writer.WritePaddedStringLine(etymologyHeader, 1)

		for _, etymology := range entry.Etymologies {
			writer.WriteStringLine(etymology)
		}

		writer.WriteNewLine()
	}
}

func printThesaurusValues(writer *defineio.PanicWriter, values source.ThesaurusValues) {
	if 0 < len(values.Synonyms) {
		writer.WritePaddedStringLine(synonymHeader, 1)

		writer.WriteStringLine(strings.Join(values.Synonyms, " ; "))

		writer.WriteNewLine()
	}

	if 0 < len(values.Antonyms) {
		writer.WritePaddedStringLine(antonymHeader, 1)

		writer.WriteStringLine(strings.Join(values.Antonyms, " ; "))

		writer.WriteNewLine()
	}
}

func getHeader(result source.DictionaryResult) string {
	firstEntry := result.Entries[0]
	header := firstEntry.Word

	if len(firstEntry.Pronunciations) > 0 {
		header = fmt.Sprintf("%s  %s", header, firstEntry.Pronunciations)
	}

	return header
}

func getEntryHeader(resultHeader string, lastEntryHeader string, lastWord string, entry source.DictionaryEntry) string {
	var header string

	if len(entry.Pronunciations) > 0 {
		header = fmt.Sprintf("%s  %s", entry.Word, entry.Pronunciations)
	} else if entry.Word != lastWord {
		header = entry.Word
	}

	if header == resultHeader || header == lastEntryHeader {
		return ""
	}

	return header
}

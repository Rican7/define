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

// PrintResult prints a source.Result.
func (p *ResultPrinter) PrintResult(result source.Result) {
	p.out.IndentWrites(func(writer *defineio.PanicWriter) {
		writer.WritePaddedStringLine(getHeader(result), 1)

		for _, entry := range result.Entries() {
			if entryHeader := getEntryHeader(result, entry); entryHeader != "" {
				writer.WriteNewLine()
				writer.WriteNewLine()
				writer.WriteStringLine(entryHeader)
			}

			writer.IndentWrites(func(writer *defineio.PanicWriter) {
				printEntry(writer, entry)
			})
		}

		writer.WriteNewLine()
	})
}

func printEntry(writer *defineio.PanicWriter, entry source.DictionaryEntry) {
	if wordEntry, isWordEntry := entry.(source.WordEntry); isWordEntry && wordEntry.Category() != "" {
		writer.WritePaddedStringLine(fmt.Sprintf("(%s)", wordEntry.Category()), 1)
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

		writer.IndentWritesBy(uint(len(prefix)), func(writer *defineio.PanicWriter) {
			for _, examples := range sense.Examples() {
				writer.WriteStringLine(fmt.Sprintf("%q", examples))
			}

			for _, notes := range sense.Notes() {
				writer.WriteStringLine(fmt.Sprintf("[%s]", notes))
			}
		})

		writer.IndentWrites(func(writer *defineio.PanicWriter) {
			for _, subSense := range sense.Subsenses() {
				prefix := " - "

				for _, definition := range subSense.Definitions() {
					writer.WriteStringLine(prefix + definition)
				}

				writer.IndentWritesBy(uint(len(prefix)), func(writer *defineio.PanicWriter) {
					if len(subSense.Examples()) > 0 {
						writer.WriteStringLine(fmt.Sprintf("%q", subSense.Examples()[0]))
					}
				})
			}
		})
	}

	if etymologyEntry, ok := entry.(source.EtymologyEntry); ok {
		printEtymologyEntry(writer, etymologyEntry)
	}

	if thesaurusEntry, ok := entry.(source.ThesaurusEntry); ok {
		printThesaurusEntry(writer, thesaurusEntry)
	}
}

func printEtymologyEntry(writer *defineio.PanicWriter, entry source.EtymologyEntry) {
	if 0 < len(entry.Etymologies()) {
		writer.WritePaddedStringLine(etymologyHeader, 1)

		for _, etymology := range entry.Etymologies() {
			writer.WriteStringLine(etymology)
		}

		writer.WriteNewLine()
	}
}

func printThesaurusEntry(writer *defineio.PanicWriter, entry source.ThesaurusEntry) {
	if 0 < len(entry.Synonyms()) {
		writer.WritePaddedStringLine(synonymHeader, 1)

		writer.WriteStringLine(strings.Join(entry.Synonyms(), " ; "))

		writer.WriteNewLine()
	}

	if 0 < len(entry.Antonyms()) {
		writer.WritePaddedStringLine(antonymHeader, 1)

		writer.WriteStringLine(strings.Join(entry.Antonyms(), " ; "))

		writer.WriteNewLine()
	}
}

func getHeader(result source.Result) string {
	header := result.Headword()

	firstEntry := result.Entries()[0]

	if firstEntry.Pronunciation() != "" || (isSameWord(result, firstEntry) && firstEntry.Pronunciation() != "") {
		header = fmt.Sprintf("%s  /%s/", header, firstEntry.Pronunciation())
	}

	return header
}

func getEntryHeader(result source.Result, entry source.DictionaryEntry) string {
	var header string

	if wordEntry, isWordEntry := entry.(source.WordEntry); isWordEntry && !isSameWord(result, entry) {
		if entry.Pronunciation() != "" {
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

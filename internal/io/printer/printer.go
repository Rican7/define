// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package printer provides types and methods to encapsulate consistent printing
// of structures.
package printer

import (
	"fmt"
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
	out             *defineio.PanicWriter
	indentationSize uint
}

// NewResultPrinter creates a new ResultPrinter.
func NewResultPrinter(out *defineio.PanicWriter, indentationSize uint) *ResultPrinter {
	return &ResultPrinter{out: out, indentationSize: indentationSize}
}

// PrintResult prints a source.Result.
func (p *ResultPrinter) PrintResult(result source.Result) {
	p.out.IndentWrites(p.indentationSize, func(writer *defineio.PanicWriter) {
		writer.WriteNewLine()
		writer.WriteStringLine(getHeader(result))
		writer.WriteNewLine()

		for _, entry := range result.Entries() {
			if entryHeader := getEntryHeader(result, entry); "" != entryHeader {
				writer.WriteNewLine()
				writer.WriteNewLine()
				writer.WriteStringLine(entryHeader)
			}

			writer.IndentWrites(p.indentationSize, func(writer *defineio.PanicWriter) {

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

					writer.IndentWrites(p.indentationSize, func(writer *defineio.PanicWriter) {
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
						writer.WriteStringLine(etymologyHeader)
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
						writer.WriteStringLine(synonymHeader)
						writer.WriteNewLine()

						writer.WriteStringLine(strings.Join(thesaurusEntry.Synonyms(), " ; "))

						writer.WriteNewLine()
					}

					if 0 < len(thesaurusEntry.Antonyms()) {
						writer.WriteNewLine()
						writer.WriteStringLine(antonymHeader)
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

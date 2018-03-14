// TODO
//
// Copyright © 2018 Trevor N. Suarez (Rican7)
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	defineio "github.com/Rican7/define/io"
	"github.com/Rican7/define/source"
	"github.com/Rican7/define/source/glosbe"
)

func main() {
	var word string

	if len(os.Args) >= 2 {
		word = os.Args[1]
	}

	result, err := glosbe.New(http.Client{}).Define(word)

	if nil != err {
		panic(err)
	}

	printResult(result, os.Stdout)
}

func printResult(result source.Result, out io.Writer) {
	writer := defineio.NewPanicWriter(out)

	const indentSize = 2

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
				}

				if thesaurusEntry, ok := entry.(source.ThesaurusEntry); ok {
					writer.WriteNewLine()

					for _, synonym := range thesaurusEntry.Synonyms() {
						writer.WriteStringLine(synonym)
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

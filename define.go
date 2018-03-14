// TODO
//
// Copyright © 2018 Trevor N. Suarez (Rican7)
package main

import (
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
	writer := defineio.PanicWriter{out}

	writer.WriteStringLine(result.Headword())
	writer.WriteNewLine()

	writer.FWriteln(len(result.Entries()))

	for _, entry := range result.Entries() {
		// TODO
		writer.WriteStringLine(entry.Pronounciation())

		for _, sense := range entry.Senses() {
			for _, definition := range sense.Definitions() {
				writer.WriteStringLine(definition)
			}
		}

		if thesaurusEntry, ok := entry.(source.ThesaurusEntry); ok {
			writer.WriteNewLine()

			for _, synonym := range thesaurusEntry.Synonyms() {
				writer.WriteStringLine(synonym)
			}
		}
	}
}

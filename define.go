// TODO
//
// Copyright © 2018 Trevor N. Suarez (Rican7)
package main

import (
	"fmt"
	"net/http"
	"os"

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

	fmt.Println(result.Headword())
	fmt.Println()

	for _, entry := range result.Entries() {
		for _, sense := range entry.Senses() {
			for _, definition := range sense.Definitions() {
				fmt.Println(definition)
			}
		}

		if thesaurusEntry, ok := entry.(source.ThesaurusEntry); ok {
			fmt.Println()

			for _, synonym := range thesaurusEntry.Synonyms() {
				fmt.Println(synonym)
			}
		}
	}
}

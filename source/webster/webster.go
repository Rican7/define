// Package webster provides a dictionary source via the Webster Dictionaries API
//
// Copyright © 2018 Trevor N. Suarez (Rican7)
package webster

import (
	"encoding/xml"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Rican7/define/source"
)

const (
	// baseURLString is the base URL for all Webster API interactions
	baseURLString = "http://www.dictionaryapi.com/api/v1/"

	entriesURLString = baseURLString + "references/collegiate/xml/"

	httpRequestAcceptHeaderName     = "Accept"
	httpRequestAppKeyQueryParamName = "key"

	xmlMIMEType = "application/xml"

	senseTagName        = "sn"
	senseDividerTagName = "sd"
	definingTextTagName = "dt"
	calledAlsoTagName   = "ca"

	senseDividerPrefix       = "; "
	definingTextPrefix       = ":"
	exampleWrapCharacter     = "'"
	authorPrefix             = "- "
	crossReferenceJoinString = ", "
)

// apiURL is the URL instance used for Webster API calls
var apiURL *url.URL

// stringCleaner is used to clean the strings returned from the API
// TODO: Replace with an HTML entities cleaner?
var stringCleaner *strings.Replacer

// tagsToClean is a list of HTML tags (start and end) to clean/remove from the
// strings returned from the API
var tagsToClean = []string{
	"<it>",
	"</it>",
	"<fw>",
	"</fw>",
	"<d_link>",
	"</d_link>",
	"<dx>",
	"</dx>",
	"<dxt>",
	"</dxt>",
	"<sx>",
	"</sx>",
	"<ca>",
	"</ca>",
	"<cat>",
	"</cat>",
	"<g>",
	"</g>",
	"<un>",
	"</un>",
	"<vi>",
	"</vi>",
	"<aq>",
	"</aq>",
	"<sxn>",
	"</sxn>",
	"<dxn>",
	"</dxn>",
}

// etymologyMetaStripperRegex is a regular expression for stripping meta from
// etymology entries
var etymologyMetaStripperRegex = regexp.MustCompile("<ma>.*?</ma>")

// api contains a configured HTTP client for Webster API operations
type api struct {
	httpClient *http.Client
	appKey     string
}

// apiResult defines the data structure for Webster API results
type apiResult struct {
	Entries []struct {
		ID                   string                   `xml:"id,attr"`
		Word                 string                   `xml:"ew"`
		Pronunciation        string                   `xml:"pr"`
		LexicalCategory      string                   `xml:"fl"`
		Etymologies          []cleanableString        `xml:"et"`
		DefinitionContainers []apiDefinitionContainer `xml:"def"`
	} `xml:"entry"`
}

// apiDefinitionContainer defines the data structure for Oxford API definitions
type apiDefinitionContainer struct {
	Raw           string            `xml:",innerxml"`
	Date          string            `xml:"date"`
	DefiningTexts []cleanableString `xml:"dt"`

	senses []apiSense
}

// apiDefiningText defines the data structure for defining texts
type apiDefiningText struct {
	Raw             string       `xml:",innerxml"`
	Stripped        string       `xml:",chardata"`
	CrossReferences []string     `xml:"sx"`
	Examples        []apiExample `xml:"vi"`
	UsageNotes      []struct {
		Note     string       `xml:",chardata"`
		Examples []apiExample `xml:"vi"`
	} `xml:"un"`

	cleaned   string
	formatted string
}

// apiExample defines the data structure for examples
type apiExample struct {
	Raw      string `xml:",innerxml"`
	Stripped string `xml:",chardata"`
	Author   string `xml:"aq"`

	cleaned   string
	formatted string
}

// apiSense defines the data structure for Oxford API senses
type apiSense struct {
	Definitions []string
	Examples    []string
	Notes       []string

	Subsenses []apiSense
}

// sensePosition defines the data structure for sense positions
type sensePosition struct {
	Position    string `xml:",chardata"`
	SubPosition string `xml:"snp"`
}

// cleanableString defines the data structure for cleanable XML strings
type cleanableString struct {
	Raw      string `xml:",innerxml"`
	Stripped string `xml:",chardata"`

	cleaned string
}

// websterEntry contains the entry types for this API
type websterEntry struct {
	source.WordEntryValue
	source.DictionaryEntryValue
	source.EtymologyEntryValue
}

// Initialize the package
func init() {
	var err error

	apiURL, err = url.Parse(baseURLString)

	if nil != err {
		panic(err)
	}

	stringCleanerPairs := make([]string, len(tagsToClean)*2)

	for _, tag := range tagsToClean {
		stringCleanerPairs = append(stringCleanerPairs, []string{tag, ""}...)
	}

	stringCleaner = strings.NewReplacer(stringCleanerPairs...)
}

// New returns a new Webster API dictionary source
func New(httpClient http.Client, appKey string) source.Source {
	return &api{&httpClient, appKey}
}

// Define takes a word string and returns a dictionary source.Result
func (g *api) Define(word string) (source.Result, error) {
	// Prepare our URL
	requestURL, err := url.Parse(entriesURLString + word)
	queryParams := apiURL.Query()
	queryParams.Set(httpRequestAppKeyQueryParamName, g.appKey)
	requestURL.RawQuery = queryParams.Encode()

	if nil != err {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)

	if nil != err {
		return nil, err
	}

	httpRequest.Header.Set(httpRequestAcceptHeaderName, xmlMIMEType)

	httpResponse, err := g.httpClient.Do(httpRequest)

	if nil != err {
		return nil, err
	}

	defer httpResponse.Body.Close()

	if err = source.ValidateHTTPResponse(httpResponse); nil != err {
		return nil, err
	}

	body, err := ioutil.ReadAll(httpResponse.Body)

	if nil != err {
		return nil, err
	}

	var result apiResult

	if err = xml.Unmarshal(body, &result); nil != err {
		return nil, err
	}

	if len(result.Entries) < 1 {
		return nil, &source.EmptyResultError{word}
	}

	return source.ValidateAndReturnResult(result.toResult())
}

// UnmarshalXML customizes the way we can unmarshal our API definitions value
func (s *apiDefinitionContainer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error

	// Alias our type so that we can unmarshal as usual
	type sense apiDefinitionContainer

	err = d.DecodeElement((*sense)(s), &start)

	// Create a new decoder for our RAW parts
	subDecoder := xml.NewDecoder(strings.NewReader(s.Raw))

	senses := []*apiSense{}
	subsenses := []*apiSense{}
	senseIndex := 0
	var currentSense *apiSense
	isDefinitionContinuation := false

	for token, err := subDecoder.Token(); nil != token || nil == err; token, err = subDecoder.Token() {
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case senseTagName:
				sp := &sensePosition{}
				err = subDecoder.DecodeElement(&sp, &t)

				currentSense = &apiSense{}

				// If the position is a number, then its a top-level sense
				if _, err := strconv.Atoi(sp.Position); err == nil {
					if len(subsenses) > 0 {
						senses[senseIndex].Subsenses = make([]apiSense, len(subsenses))
						for i, subsense := range subsenses {
							senses[senseIndex].Subsenses[i] = *subsense
						}

						// Reset our subsenses
						subsenses = make([]*apiSense, 0)
					}

					if len(senses) > 0 {
						senseIndex++
					}
					senses = append(senses, currentSense)
				} else {
					subsenses = append(subsenses, currentSense)
				}
			case senseDividerTagName:
				lastDefinitionIndex := len(currentSense.Definitions) - 1

				str := &cleanableString{}
				err = subDecoder.DecodeElement(&str, &t)

				currentSense.Definitions[lastDefinitionIndex] =
					currentSense.Definitions[lastDefinitionIndex] + senseDividerPrefix + str.cleaned

				isDefinitionContinuation = true
			case definingTextTagName:
				if len(senses) == 0 || nil == currentSense {
					currentSense = &apiSense{}
					senses = append(senses, currentSense)
				}

				dt := &apiDefiningText{}
				err = subDecoder.DecodeElement(&dt, &t)

				if !isDefinitionContinuation {
					currentSense.Definitions = append(currentSense.Definitions, dt.formatted)

					currentSense.Examples = make([]string, len(dt.Examples))
					for i, example := range dt.Examples {
						currentSense.Examples[i] = example.formatted
					}

					currentSense.Notes = make([]string, len(dt.UsageNotes))
					for i, note := range dt.UsageNotes {
						currentSense.Notes[i] = note.Note
					}
				} else {
					lastDefinitionIndex := len(currentSense.Definitions) - 1

					currentSense.Definitions[lastDefinitionIndex] =
						currentSense.Definitions[lastDefinitionIndex] + " " + dt.formatted

					isDefinitionContinuation = false
				}
			}

			if nil != err {
				return err
			}
		}
	}

	s.senses = make([]apiSense, len(senses))
	for i, sense := range senses {
		s.senses[i] = *sense
	}

	return err
}

// UnmarshalXML customizes the way we can unmarshal our API defining texts value
func (dt *apiDefiningText) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Alias our type so that we can unmarshal as usual
	type defText apiDefiningText

	// Unmarshal our wrapped value before cleaning
	err := d.DecodeElement((*defText)(dt), &start)

	// Initialize our cleaned string
	cleanStr := &cleanableString{}
	xml.Unmarshal(wrapRawXML(dt.Raw), cleanStr)
	dt.cleaned = cleanStr.cleaned

	dt.cleaned = strings.TrimLeft(dt.cleaned, definingTextPrefix)

	// Our formatted version will start as just the cleaned version
	dt.formatted = dt.cleaned

	// Clean our cross references
	for i, crossReference := range dt.CrossReferences {
		crossReference = stringCleaner.Replace(crossReference)
		crossReference = strings.TrimSpace(crossReference)
		crossReference = strings.TrimLeft(crossReference, definingTextPrefix)

		dt.CrossReferences[i] = crossReference
	}

	for i, usageNote := range dt.UsageNotes {
		// Grab our examples from our usage notes
		dt.Examples = append(dt.Examples, usageNote.Examples...)

		// Clean our note
		dt.UsageNotes[i].Note = strings.TrimSpace(usageNote.Note)
	}

	// If we only have a single usage note, and the defining text starts with it
	if len(dt.UsageNotes) == 1 && strings.HasPrefix(dt.cleaned, dt.UsageNotes[0].Note) {
		// Functionally replace the defining text with the note
		dt.formatted = dt.UsageNotes[0].Note

		// Remove the note, since it would then be redundant
		dt.UsageNotes = dt.UsageNotes[:0]
	} else {
		for _, usageNote := range dt.UsageNotes {
			if strings.Contains(dt.formatted, usageNote.Note) {
				parts := strings.SplitN(dt.formatted, usageNote.Note, 2)

				// Get our start and end pieces
				strStart := strings.TrimSpace(parts[0])
				strEnd := strings.TrimSpace(parts[1])

				dt.formatted = strStart + strEnd
			}
		}
	}

	for _, example := range dt.Examples {
		if strings.Contains(dt.formatted, example.cleaned) {
			parts := strings.SplitN(dt.formatted, example.cleaned, 2)

			// Get our start and end pieces
			strStart := strings.TrimSpace(parts[0])
			strEnd := strings.TrimSpace(parts[1])

			dt.formatted = strStart + strEnd
		}
	}

	// If our cleaned string only contains our cross references
	if len(dt.CrossReferences) > 1 && strings.EqualFold(dt.formatted, strings.Join(dt.CrossReferences, " ")) {
		// Add commas, for readability
		dt.formatted = strings.Join(dt.CrossReferences, crossReferenceJoinString)
	}

	return err
}

// UnmarshalXML customizes the way we can unmarshal our API example value
func (e *apiExample) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Alias our type so that we can unmarshal as usual
	type example apiExample

	// Unmarshal our wrapped value before cleaning
	err := d.DecodeElement((*example)(e), &start)

	// Initialize our cleaned string
	cleanStr := &cleanableString{}
	xml.Unmarshal(wrapRawXML(e.Raw), cleanStr)
	e.cleaned = cleanStr.cleaned

	// Our formatted version will start as just the cleaned version
	e.formatted = e.cleaned

	// Clean our author string
	e.Author = strings.TrimSpace(e.Author)

	// If we have an author
	if "" != e.Author {
		// If the author is in the string, strip it from the original string,
		// so that we can properly append it
		if strings.Contains(e.cleaned, e.Author) {
			parts := strings.SplitN(e.cleaned, e.Author, 2)

			// Get our start and end pieces
			strStart := strings.TrimSpace(parts[0])
			strEnd := strings.TrimSpace(parts[1])

			// If we have an ending string, pad it
			if 0 < len(strEnd) {
				strEnd = " " + strEnd
			}

			e.formatted = exampleWrapCharacter + strStart + exampleWrapCharacter + " " + authorPrefix + e.Author + strEnd
		} else {
			e.formatted = exampleWrapCharacter + e.cleaned + exampleWrapCharacter + " " + authorPrefix + e.Author
		}
	}

	return err
}

// UnmarshalXML customizes the way we can unmarshal cleanable strings
func (s *cleanableString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Alias our type so that we can unmarshal as usual
	type str cleanableString

	// Unmarshal our wrapped value before cleaning
	err := d.DecodeElement((*str)(s), &start)

	// Initialize our clean string
	s.cleaned = s.Raw

	// Clean our raw string
	s.cleaned = html.UnescapeString(s.cleaned)
	s.cleaned = stringCleaner.Replace(s.cleaned)
	s.cleaned = strings.TrimSpace(s.cleaned)

	return err
}

// toResult converts the proprietary API result to a generic source.Result
func (r apiResult) toResult() source.Result {
	mainEntry := r.Entries[0]
	mainWord := mainEntry.Word

	entries := make([]interface{}, 0)

	for _, apiEntry := range r.Entries {
		if !strings.EqualFold(apiEntry.Word, mainWord) {
			continue
		}

		entry := &websterEntry{}

		entry.WordVal = apiEntry.Word
		entry.PronunciationVal = apiEntry.Pronunciation
		entry.CategoryVal = apiEntry.LexicalCategory

		entry.EtymologyVals = make([]string, len(apiEntry.Etymologies))
		for i, etymology := range apiEntry.Etymologies {
			etymology.cleaned = etymologyMetaStripperRegex.ReplaceAllString(etymology.cleaned, "")

			entry.EtymologyVals[i] = strings.TrimSpace(etymology.cleaned)
		}

		if len(apiEntry.DefinitionContainers) > 0 {
			def := apiEntry.DefinitionContainers[0]

			for _, sense := range def.senses {
				senseVal := sense.toSenseValue()

				// Only go one level deep of sub-senses
				for _, subSense := range sense.Subsenses {
					senseVal.SubsenseVals = append(senseVal.SubsenseVals, subSense.toSenseValue())
				}

				entry.SenseVals = append(entry.SenseVals, senseVal)
			}
		}

		entries = append(entries, entry)
	}

	return source.ResultValue{
		Head:      mainWord,
		Lang:      "en", // TODO
		EntryVals: entries,
	}
}

// toSenseValue converts the proprietary API sense to a source.SenseValue
func (s apiSense) toSenseValue() source.SenseValue {
	return source.SenseValue{
		DefinitionVals: s.Definitions,
		ExampleVals:    s.Examples,
		NoteVals:       s.Notes,
	}
}

// wrapRawXML wraps a raw XML string in arbitrary container elements
func wrapRawXML(raw string) []byte {
	return []byte("<x>" + raw + "</x>")
}

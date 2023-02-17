// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package webster provides a dictionary source via the Webster Dictionaries API
package webster

import (
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Rican7/define/source"
	"github.com/microcosm-cc/bluemonday"
)

// Name defines the name of the source
const Name = "Merriam-Webster's Dictionary API"

const (
	// baseURLString is the base URL for all Webster API interactions
	baseURLString = "http://www.dictionaryapi.com/api/v1/"

	entriesURLString = baseURLString + "references/collegiate/xml/"

	httpRequestAcceptHeaderName     = "Accept"
	httpRequestAppKeyQueryParamName = "key"

	xmlMIMEType     = "application/xml"
	xmlTextMIMEType = "text/xml"
	xmlBaseMIMEType = "xml"

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

// validMIMETypes is the list of valid response MIME types
var validMIMETypes = []string{xmlMIMEType, xmlTextMIMEType, xmlBaseMIMEType}

// htmlCleaner is used to clean the strings returned from the API
var htmlCleaner = bluemonday.StrictPolicy()

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

// Initialize the package
func init() {
	var err error

	apiURL, err = url.Parse(baseURLString)

	if err != nil {
		panic(err)
	}
}

// New returns a new Webster API dictionary source
func New(httpClient http.Client, appKey string) source.Source {
	return &api{&httpClient, appKey}
}

// Name returns the printable, human-readable name of the source.
func (g *api) Name() string {
	return Name
}

// Define takes a word string and returns a list of dictionary results, and
// an error if any occurred.
func (g *api) Define(word string) ([]source.DictionaryResult, error) {
	// Prepare our URL
	requestURL, err := url.Parse(entriesURLString + word)
	queryParams := apiURL.Query()
	queryParams.Set(httpRequestAppKeyQueryParamName, g.appKey)
	requestURL.RawQuery = queryParams.Encode()

	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(httpRequestAcceptHeaderName, xmlMIMEType)
	httpRequest.Header.Add(httpRequestAcceptHeaderName, xmlTextMIMEType)
	httpRequest.Header.Add(httpRequestAcceptHeaderName, xmlBaseMIMEType)

	httpResponse, err := g.httpClient.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	if err = source.ValidateHTTPResponse(httpResponse, validMIMETypes, nil); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		return nil, err
	}

	var result apiResult

	if err = xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Entries) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return source.ValidateAndReturnDictionaryResults(word, result.toResults())
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

	for token, err := subDecoder.Token(); token != nil || err == nil; token, err = subDecoder.Token() {
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
				if len(senses) == 0 || currentSense == nil {
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

			if err != nil {
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
		crossReference = htmlCleaner.Sanitize(crossReference)
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
	if e.Author != "" {
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
	s.cleaned = htmlCleaner.Sanitize(s.cleaned)
	s.cleaned = html.UnescapeString(s.cleaned)
	s.cleaned = strings.TrimSpace(s.cleaned)

	return err
}

// toResult converts the API response to the results that a source expects to
// return.
func (r *apiResult) toResults() []source.DictionaryResult {
	mainEntry := r.Entries[0]
	mainWord := mainEntry.Word

	sourceEntries := make([]source.DictionaryEntry, 0, len(r.Entries))

	for _, apiEntry := range r.Entries {
		if !strings.EqualFold(apiEntry.Word, mainWord) {
			continue
		}

		entry := source.DictionaryEntry{}

		entry.Word = apiEntry.Word
		entry.LexicalCategory = apiEntry.LexicalCategory

		if apiEntry.Pronunciation != "" {
			entry.Pronunciations = []string{apiEntry.Pronunciation}
		}

		entry.Etymologies = make([]string, 0, len(apiEntry.Etymologies))
		for _, etymology := range apiEntry.Etymologies {
			etymology.cleaned = etymologyMetaStripperRegex.ReplaceAllString(etymology.cleaned, "")

			entry.Etymologies = append(entry.Etymologies, strings.TrimSpace(etymology.cleaned))
		}

		if len(apiEntry.DefinitionContainers) > 0 {
			def := apiEntry.DefinitionContainers[0]

			for _, sense := range def.senses {
				sourceSense := sense.toSense()

				// Only go one level deep of sub-senses
				for _, subSense := range sense.Subsenses {
					sourceSense.SubSenses = append(sourceSense.SubSenses, subSense.toSense())
				}

				entry.Senses = append(entry.Senses, sourceSense)
			}
		}

		sourceEntries = append(sourceEntries, entry)
	}

	return []source.DictionaryResult{
		{
			Language: "en", // TODO
			Entries:  sourceEntries,
		},
	}
}

// toSense converts the API sense to a source.Sense
func (s *apiSense) toSense() source.Sense {
	return source.Sense{
		Definitions: s.Definitions,
		Examples:    s.Examples,
		Notes:       s.Notes,
	}
}

// wrapRawXML wraps a raw XML string in arbitrary container elements
func wrapRawXML(raw string) []byte {
	return []byte("<x>" + raw + "</x>")
}

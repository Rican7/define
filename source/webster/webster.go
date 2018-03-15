// Package webster provides a dictionary source via the Webster Dictionaries API
//
// Copyright © 2018 Trevor N. Suarez (Rican7)
package webster

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Rican7/define/source"
)

const (
	// baseURLString is the base URL for all Webster API interactions
	baseURLString = "http://www.dictionaryapi.com/api/v1/"

	entriesURLString = baseURLString + "references/collegiate/xml/"

	senseTagName        = "sn"
	senseDividerTagName = "sd"
	definingTextTagName = "dt"
	calledAlsoTagName   = "ca"

	definingTextPrefix = ":"
	calledAlsoPrefix   = "- "
	senseDividerPrefix = "; "
)

// apiURL is the URL instance used for Webster API calls
var apiURL *url.URL

// api is a struct containing a configured HTTP client for Webster API operations
type api struct {
	httpClient *http.Client
	appKey     string
}

// apiResult is a struct that defines the data structure for Webster API results
type apiResult struct {
	Entries []struct {
		ID              string `xml:"id,attr"`
		Word            string `xml:"ew"`
		Pronunciation   string `xml:"pr"`
		LexicalCategory string `xml:"fl"`
		Etymologies     []struct {
			Raw       string `xml:",innerxml"`
			Etymology string `xml:",chardata"`
		} `xml:"et"`
		DefinitionContainers []apiDefinitionContainer `xml:"def"`
	} `xml:"entry"`
}

// apiDefinitionContainer is a struct that defines the data structure for
// Oxford API definitions
type apiDefinitionContainer struct {
	Raw           string `xml:",innerxml"`
	Date          string `xml:"date"`
	DefiningTexts []struct {
		Raw          string `xml:",innerxml"`
		DefiningText string `xml:",chardata"`
	} `xml:"dt"`

	Senses []apiSense // TODO make this private
}

// apiSense is a struct that defines the data structure for Oxford API senses
type apiSense struct {
	Definitions []string

	Subsenses []apiSense
}

// sensePosition is a struct that defines the data structure for sense positions
type sensePosition struct {
	Position    string `xml:",chardata"`
	SubPosition string `xml:"snp"`
}

// sensePosition is a struct that defines the data structure for defining texts
type definingText struct {
	Raw             string   `xml:",innerxml"`
	DefiningText    string   `xml:",chardata"`
	CrossReferences []string `xml:"sx"`

	cleaned string
}

// websterEntry is a struct that contains the entry types for this API
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
	queryParams.Set("key", g.appKey)
	requestURL.RawQuery = queryParams.Encode()

	if nil != err {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)

	if nil != err {
		return nil, err
	}

	httpRequest.Header.Set("Accept", "application/xml")

	httpResponse, err := g.httpClient.Do(httpRequest)

	if nil != err {
		return nil, err
	}

	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)

	if nil != err {
		return nil, err
	}

	var result apiResult
	err = xml.Unmarshal(body, &result)

	return result.toResult(), err
}

// UnmarshalXML customizes the way we can unmarshal our API result value
func (r *apiResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Alias our type so that we can unmarshal as usual
	type result apiResult

	// Unmarshal our wrapped value before cleaning
	err := d.DecodeElement((*result)(r), &start)

	// TODO: Replace with an HTML entities cleaner?
	strCleaner := strings.NewReplacer(
		"<it>", "",
		"</it>", "",
	)

	for _, entry := range r.Entries {
		for i, etymology := range entry.Etymologies {
			// Clean the strings
			etymology.Etymology = strCleaner.Replace(etymology.Raw)

			// Store the modified value
			entry.Etymologies[i] = etymology
		}
	}

	return err
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

				dt := &definingText{}
				err = subDecoder.DecodeElement(&dt, &t)

				currentSense.Definitions[lastDefinitionIndex] =
					currentSense.Definitions[lastDefinitionIndex] + senseDividerPrefix + dt.cleaned

				isDefinitionContinuation = true
			case definingTextTagName:
				if len(senses) == 0 || nil == currentSense {
					currentSense = &apiSense{}
					senses = append(senses, currentSense)
				}

				dt := &definingText{}
				err = subDecoder.DecodeElement(&dt, &t)

				if !isDefinitionContinuation {
					currentSense.Definitions = append(currentSense.Definitions, dt.cleaned)
				} else {
					lastDefinitionIndex := len(currentSense.Definitions) - 1

					currentSense.Definitions[lastDefinitionIndex] =
						currentSense.Definitions[lastDefinitionIndex] + " " + dt.cleaned

					isDefinitionContinuation = false
				}
			}

			if nil != err {
				return err
			}
		}
	}

	s.Senses = make([]apiSense, len(senses))
	for i, sense := range senses {
		s.Senses[i] = *sense
	}

	return err
}

// UnmarshalXML customizes the way we can unmarshal our API defining texts value
func (dt *definingText) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Alias our type so that we can unmarshal as usual
	type defText definingText

	// Unmarshal our wrapped value before cleaning
	err := d.DecodeElement((*defText)(dt), &start)

	// TODO: Replace with an HTML entities cleaner?
	strCleaner := strings.NewReplacer(
		"<it>", "",
		"</it>", "",
		"<d_link>", "",
		"</d_link>", "",
		"<dx>", "",
		"</dx>", "",
		"<dxt>", "",
		"</dxt>", "",
		"<sx>", "",
		"</sx>", "",
		"<ca>", "",
		"</ca>", "",
		"<cat>", "",
		"</cat>", "",
		"<g>", "",
		"</g>", "",
	)

	// Clean our raw string
	dt.cleaned = strCleaner.Replace(dt.Raw)
	dt.cleaned = strings.TrimSpace(dt.cleaned)
	dt.cleaned = strings.TrimLeft(dt.cleaned, definingTextPrefix)

	// Clean our cross references
	for i, crossReference := range dt.CrossReferences {
		crossReference = strCleaner.Replace(crossReference)
		crossReference = strings.TrimSpace(crossReference)
		crossReference = strings.TrimLeft(crossReference, definingTextPrefix)

		dt.CrossReferences[i] = crossReference
	}

	// If our cleaned string only contains our cross references
	if len(dt.CrossReferences) > 1 && strings.EqualFold(dt.cleaned, strings.Join(dt.CrossReferences, " ")) {
		// Add commas, for readability
		dt.cleaned = strings.Join(dt.CrossReferences, ", ")
	}

	return err
}

// toResult converts the proprietary API result to a generic source.Result
func (r apiResult) toResult() source.Result {
	mainEntry := r.Entries[0] // TODO: Handle empty entries as an error in the `Define` method
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
			entry.EtymologyVals[i] = etymology.Etymology
		}

		if len(apiEntry.DefinitionContainers) > 0 {
			def := apiEntry.DefinitionContainers[0]

			for _, sense := range def.Senses {
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
	}
}

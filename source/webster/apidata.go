package webster

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Rican7/define/source"
)

const (
	// See https://www.dictionaryapi.com/products/json#sec-2
	arrayDataTagText                = "text"
	arrayDataTagSense               = "sense"
	arrayDataTagBindingSubstitute   = "bs"
	arrayDataTagVerbalIllustrations = "vis"

	// See https://www.dictionaryapi.com/products/json#sec-2
	objectDataTagText               = "t"
	objectDataTagSense              = "sense"
	objectDataTagSenseNumber        = "sn"
	objectDataTagDefiningText       = "dt"
	objectDataTagAttributionOfQuote = "aq"
	objectDataTagAuthor             = "auth"
)

var (
	// regexpWebsterTokens is a regular exprssion for matching Webster API
	// text tokens.
	//
	// Webster API text may contain "tokens", as documented:
	//
	// See https://www.dictionaryapi.com/products/json#sec-2.tokens
	regexpWebsterTokens = regexp.MustCompile(`{.*?(?:\|(.*?)(?:\|.*?\|?)?)?}`)

	// regexpWebsterSenseNumber is a regular expression for matching Webster API
	// sense numbers.
	//
	// Webster API sense numbers may be numerals, lower-case letters, or
	// parenthesized numerals, as documented:
	//
	// See https://www.dictionaryapi.com/products/json#sec-2.sn
	regexpWebsterSenseNumber = regexp.MustCompile(`(\d+)? ?(\w+)? ?(\(\d+\))?`)
)

// apiRawResponse defines the data structure for a raw Webster API response
type apiRawResponse []any

// apiResponse defines the data structure for a Webster API response
type apiResponse[T apiResponseItem] []T

// apiResponseItem defines a type constraint for Webster API response items
type apiResponseItem interface {
	apiSearchResult | apiDefinitionResult
}

// apiSearchResults defines the data for Webster API search results
type apiSearchResults []apiSearchResult

// apiDefinitionResults defines the data for Webster API definition results
type apiDefinitionResults []apiDefinitionResult

// apiSearchResult defines the data for a Webster API search result
type apiSearchResult string

// apiDefinitionResult defines the data structure for a Webster API definition
// result
type apiDefinitionResult struct {
	Meta apiDefinitionMeta         `json:"meta"`
	Hom  int                       `json:"hom"`
	Hwi  apiDefinitionHeadwordInfo `json:"hwi"`
	Fl   string                    `json:"fl"`
	Ins  []struct {
		If string `json:"if"`
	} `json:"ins"`
	Def  []apiDefinitionSectionEntry `json:"def"`
	Uros []struct {
		Ure string `json:"ure"`
		Fl  string `json:"fl"`
	} `json:"uros"`
	Syns []struct {
		Pl string  `json:"pl"`
		Pt [][]any `json:"pt"`
	} `json:"syns"`
	Et     [][]string `json:"et"`
	Date   string     `json:"date"`
	LdLink struct {
		LinkHw string `json:"link_hw"`
		LinkFl string `json:"link_fl"`
	} `json:"ld_link"`
	Shortdef []string `json:"shortdef"`
}

// apiDefinitionMeta defines the data structure for Webster API definition meta
type apiDefinitionMeta struct {
	ID        string   `json:"id"`
	UUID      string   `json:"uuid"`
	Sort      string   `json:"sort"`
	Src       string   `json:"src"`
	Section   string   `json:"section"`
	Stems     []string `json:"stems"`
	Offensive bool     `json:"offensive"`
}

// apiDefinitionMeta defines the data structure for Webster API definition
// headword information
type apiDefinitionHeadwordInfo struct {
	Hw  string `json:"hw"`
	Prs []struct {
		Mw    string `json:"mw"`
		Sound struct {
			Audio string `json:"audio"`
			Ref   string `json:"ref"`
			Stat  string `json:"stat"`
		} `json:"sound"`
	} `json:"prs"`
}

// apiDefinitionSectionEntry defines the data structure for Webster API
// definition section entries
type apiDefinitionSectionEntry struct {
	Vd   string           `json:"vd"`
	Sseq apiSenseSequence `json:"sseq"`
}

// apiSenseSequence defines the data structure for a Webster API sense sequence
type apiSenseSequence []apiSense

// apiSense defines the data structure for a Webster API sense
type apiSense [][]any

// apiSenseData defines the data structure for a Webster API sense data
type apiSenseData map[string]any

// apiSenseNumber defines the data structure for a Webster API sense number
type apiSenseNumber struct {
	number int
	letter string
	sub    string
}

// UnmarshalJSON satisfies the encoding/json.Unmarshaler interface
func (r *apiRawResponse) UnmarshalJSON(data []byte) error {
	var rawSlice []json.RawMessage

	if err := json.Unmarshal(data, &rawSlice); err != nil {
		return err
	}

	if len(rawSlice) < 1 || len(rawSlice[0]) < 1 {
		return nil
	}

	var newResponse apiRawResponse
	var err error

	// Inspect the first byte of the first item
	switch rawSlice[0][0] {
	case '"':
		var response apiResponse[apiSearchResult]
		newResponse, err = replaceData(data, response)
	case '{':
		var response apiResponse[apiDefinitionResult]
		newResponse, err = replaceData(data, response)
	}

	if err != nil {
		return err
	}
	*r = newResponse

	return nil
}

// replaceData takes raw JSON bytes and a new response type and returns a
// apiRawResponse with the new response data after unmarshalling.
func replaceData[T apiResponseItem](data []byte, response apiResponse[T]) (apiRawResponse, error) {
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	newResponse := make(apiRawResponse, len(response))
	for i, item := range response {
		newResponse[i] = item
	}

	return newResponse, nil
}

func apiResponseFromRaw[T apiResponseItem](raw apiRawResponse) apiResponse[T] {
	response := make(apiResponse[T], len(raw))

	if len(raw) < 1 {
		return response
	}

	for i, item := range raw {
		response[i] = item.(T)
	}

	return response
}

// toResult converts the API response to the results that a source expects to
// return.
func (r apiDefinitionResults) toResults() []source.DictionaryResult {
	mainEntry := r[0]
	mainWord := cleanHeadword(mainEntry.Hwi.Hw)

	sourceEntries := make([]source.DictionaryEntry, 0, len(r))

	for _, apiEntry := range r {
		headword := cleanHeadword(apiEntry.Hwi.Hw)

		if !source.EqualFoldPlain(headword, mainWord) {
			continue
		}

		sourceEntry := source.DictionaryEntry{}

		sourceEntry.Word = headword
		sourceEntry.LexicalCategory = apiEntry.Fl

		sourceEntry.Pronunciations = make([]string, 0, len(apiEntry.Hwi.Prs))
		for _, pronunciation := range apiEntry.Hwi.Prs {
			sourceEntry.Pronunciations = append(sourceEntry.Pronunciations, pronunciation.Mw)
		}

		// We filter etymologies and potentially add them later so the capacity
		// can't be reasonably known here.
		sourceEntry.Etymologies = make([]string, 0)
		for _, etymology := range apiEntry.Et {
			// Webster API etymologies are returned in prefixed arrays.
			// See https://www.dictionaryapi.com/products/json#sec-2.et
			if len(etymology) < 2 || etymology[0] != arrayDataTagText {
				continue
			}

			etymologyText := cleanTextOfTokens(etymology[1])

			sourceEntry.Etymologies = append(sourceEntry.Etymologies, etymologyText)
		}

		for _, def := range apiEntry.Def {
			sourceEntry.Senses = append(sourceEntry.Senses, def.Sseq.toSenses()...)
		}

		sourceEntries = append(sourceEntries, sourceEntry)
	}

	return []source.DictionaryResult{
		{
			Language: "en", // TODO
			Entries:  sourceEntries,
		},
	}
}

// toSenses converts the API sense sequence to a list of source.Sense
func (s apiSenseSequence) toSenses() []source.Sense {
	senses := make([]source.Sense, 0)

	for _, apiSense := range s {
		var lastSenseNumber *apiSenseNumber

		for _, apiSenseContainer := range apiSense {
			// Webster API senses are returned in prefixed arrays.
			// See https://www.dictionaryapi.com/products/json#sec-2.sense
			if len(apiSenseContainer) < 2 {
				continue
			}

			var senseData apiSenseData

			switch apiSenseContainer[0] {
			case arrayDataTagSense:
				senseData = apiSenseData(apiSenseContainer[1].(map[string]any))
			case arrayDataTagBindingSubstitute:
				// See https://www.dictionaryapi.com/products/json#sec-2.bs
				bindingSubstitute := apiSenseContainer[1].(map[string]any)
				senseData = apiSenseData(bindingSubstitute[objectDataTagSense].(map[string]any))
			default:
				continue
			}

			senseNumber := parseSenseNumber(senseData[objectDataTagSenseNumber])

			sourceSense := senseData.toSense()

			if lastSenseNumber == nil || (senseNumber != nil && lastSenseNumber.number < senseNumber.number) {
				// The sense is a new sense
				senses = append(senses, sourceSense)
			} else {
				// The sense is a sub-sense
				lastSense := &(senses[len(senses)-1])
				lastSense.SubSenses = append(lastSense.SubSenses, sourceSense)
			}

			lastSenseNumber = senseNumber
		}
	}

	return senses
}

// toSense converts the API sense data to a source.Sense
func (d apiSenseData) toSense() source.Sense {
	definitions := make([]string, 0)
	examples := make([]string, 0)

	senseDefinitions := d[objectDataTagDefiningText].([]any)

	for _, defParts := range senseDefinitions {
		definition := defParts.([]any)

		// Webster API definition parts are returned in prefixed arrays.
		// See https://www.dictionaryapi.com/products/json#sec-2.dt
		if len(definition) < 2 {
			continue
		}

		switch definition[0] {
		case arrayDataTagText:
			definitionText := cleanTextOfTokens(definition[1].(string))

			definitions = append(definitions, definitionText)
		case arrayDataTagVerbalIllustrations:
			exampleTextObjects := definition[1].([]any)

			for _, exampleTextObject := range exampleTextObjects {
				exampleTextMap := exampleTextObject.(map[string]any)
				exampleText := cleanTextOfTokens(exampleTextMap[objectDataTagText].(string))

				if exampleTextMap[objectDataTagAttributionOfQuote] != nil {
					exampleAttribution := exampleTextMap[objectDataTagAttributionOfQuote].(map[string]any)

					// TODO: Handle attributions more generally... especially
					// presentation wise...
					if exampleAttribution[objectDataTagAuthor] != nil {
						exampleText = fmt.Sprintf("'%s' - %s", exampleText, exampleAttribution[objectDataTagAuthor])
					}
				}

				examples = append(examples, exampleText)
			}
		}
	}

	return source.Sense{
		Definitions: definitions,
		Examples:    examples,
	}
}

func cleanHeadword(headword string) string {
	return strings.ReplaceAll(headword, "*", "")
}

func cleanTextOfTokens(text string) string {
	return regexpWebsterTokens.ReplaceAllString(text, "$1")
}

func parseSenseNumber(rawSenseNumber any) *apiSenseNumber {
	if rawSenseNumber == nil {
		return nil
	}

	parsed := regexpWebsterSenseNumber.FindStringSubmatch(rawSenseNumber.(string))

	var main int
	if parsedMain, err := strconv.Atoi(parsed[1]); err == nil {
		main = parsedMain
	}

	return &apiSenseNumber{
		number: main,
		letter: parsed[2],
		sub:    parsed[3],
	}
}

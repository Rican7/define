// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package oxford provides a dictionary source via the Oxford Dictionaries API
package oxford

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Rican7/define/source"
)

// Name defines the name of the source
const Name = "Oxford Dictionaries API V2"

const (
	// baseURLString is the base URL for all Oxford API interactions
	baseURLString = "https://od-api.oxforddictionaries.com/api/v2/"

	entriesURLString             = baseURLString + "entries/"
	defaultEntriesSourceLanguage = "en-us"

	httpRequestAcceptHeaderName = "Accept"
	httpRequestAppIDHeaderName  = "app_id"
	httpRequestAppKeyHeaderName = "app_key"

	jsonMIMEType = "application/json"

	phoneticNotationIPAIdentifier = "IPA"
)

// apiURL is the URL instance used for Oxford API calls
var apiURL *url.URL

// validMIMETypes is the list of valid response MIME types
var validMIMETypes = []string{jsonMIMEType}

// api is a struct containing a configured HTTP client for Oxford API operations
type api struct {
	httpClient *http.Client
	appID      string
	appKey     string
}

// apiResult is a struct that defines the data structure for Oxford API results
type apiResult struct {
	Metadata struct {
	} `json:"metadata"`
	Results []struct {
		ID             string `json:"id"`
		Language       string `json:"language"`
		LexicalEntries []struct {
			Compounds []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"compounds"`
			DerivativeOf []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"derivativeOf"`
			Derivatives []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"derivatives"`
			Entries []struct {
				CrossReferenceMarkers []string `json:"crossReferenceMarkers"`
				CrossReferences       []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"crossReferences"`
				Etymologies         []string `json:"etymologies"`
				GrammaticalFeatures []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"grammaticalFeatures"`
				HomographNumber string `json:"homographNumber"`
				Inflections     []struct {
					Domains []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"domains"`
					GrammaticalFeatures []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
						Type string `json:"type"`
					} `json:"grammaticalFeatures"`
					InflectedForm   string `json:"inflectedForm"`
					LexicalCategory struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"lexicalCategory"`
					Pronunciations []apiPronunciation `json:"pronunciations"`
					Regions        []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"regions"`
					Registers []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"registers"`
				} `json:"inflections"`
				Notes []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"notes"`
				Pronunciations []apiPronunciation `json:"pronunciations"`
				Senses         []apiSense         `json:"senses"`
				VariantForms   []struct {
					Domains []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"domains"`
					Notes []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
						Type string `json:"type"`
					} `json:"notes"`
					Pronunciations []apiPronunciation `json:"pronunciations"`
					Regions        []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"regions"`
					Registers []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"registers"`
					Text string `json:"text"`
				} `json:"variantForms"`
			} `json:"entries"`
			GrammaticalFeatures []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"grammaticalFeatures"`
			Language        string `json:"language"`
			LexicalCategory struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"lexicalCategory"`
			Notes []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"notes"`
			PhrasalVerbs []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"phrasalVerbs"`
			Phrases []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				ID       string `json:"id"`
				Language string `json:"language"`
				Regions  []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"phrases"`
			Pronunciations []apiPronunciation `json:"pronunciations"`
			Text           string             `json:"text"`
			VariantForms   []struct {
				Domains []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"domains"`
				Notes []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"notes"`
				Pronunciations []apiPronunciation `json:"pronunciations"`
				Regions        []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"regions"`
				Registers []struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"registers"`
				Text string `json:"text"`
			} `json:"variantForms"`
		} `json:"lexicalEntries"`
		Pronunciations []apiPronunciation `json:"pronunciations"`
		Type           string             `json:"type"`
		Word           string             `json:"word"`
	} `json:"results"`
}

// apiSense is a struct that defines the data structure for Oxford API senses
type apiSense struct {
	Antonyms []struct {
		Domains []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"domains"`
		ID       string `json:"id"`
		Language string `json:"language"`
		Regions  []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"regions"`
		Registers []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"registers"`
		Text string `json:"text"`
	} `json:"antonyms"`
	Constructions []struct {
		Domains []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"domains"`
		Examples [][]string `json:"examples"`
		Notes    []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"notes"`
		Regions []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"regions"`
		Registers []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"registers"`
		Text string `json:"text"`
	} `json:"constructions"`
	CrossReferenceMarkers []string `json:"crossReferenceMarkers"`
	CrossReferences       []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"crossReferences"`
	Definitions []string `json:"definitions"`
	Domains     []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"domains"`
	Etymologies []string `json:"etymologies"`
	Examples    []struct {
		Definitions []string `json:"definitions"`
		Domains     []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"domains"`
		Notes []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"notes"`
		Regions []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"regions"`
		Registers []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"registers"`
		SenseIds []string `json:"senseIds"`
		Text     string   `json:"text"`
	} `json:"examples"`
	ID          string `json:"id"`
	Inflections []struct {
		Domains []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"domains"`
		GrammaticalFeatures []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"grammaticalFeatures"`
		InflectedForm   string `json:"inflectedForm"`
		LexicalCategory struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"lexicalCategory"`
		Pronunciations []apiPronunciation `json:"pronunciations"`
		Regions        []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"regions"`
		Registers []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"registers"`
	} `json:"inflections"`
	Notes []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"notes"`
	Pronunciations []apiPronunciation `json:"pronunciations"`
	Regions        []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"regions"`
	Registers []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"registers"`
	ShortDefinitions []string   `json:"shortDefinitions"`
	Subsenses        []apiSense `json:"subsenses"`
	Synonyms         []struct {
		Domains []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"domains"`
		ID       string `json:"id"`
		Language string `json:"language"`
		Regions  []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"regions"`
		Registers []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"registers"`
		Text string `json:"text"`
	} `json:"synonyms"`
	ThesaurusLinks []struct {
		EntryID string `json:"entry_id"`
		SenseID string `json:"sense_id"`
	} `json:"thesaurusLinks"`
	VariantForms []struct {
		Domains []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"domains"`
		Notes []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"notes"`
		Pronunciations []apiPronunciation `json:"pronunciations"`
		Regions        []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"regions"`
		Registers []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"registers"`
		Text string `json:"text"`
	} `json:"variantForms"`
}

// apiSense is a struct that defines the data structure for Oxford API pronunciation
type apiPronunciation struct {
	AudioFile        string   `json:"audioFile"`
	Dialects         []string `json:"dialects"`
	PhoneticNotation string   `json:"phoneticNotation"`
	PhoneticSpelling string   `json:"phoneticSpelling"`
	Regions          []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"regions"`
	Registers []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"registers"`
}

// oxfordEntry is a struct that contains the entry types for this API
type oxfordEntry struct {
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

// New returns a new Oxford API dictionary source
func New(httpClient http.Client, appID, appKey string) source.Source {
	return &api{&httpClient, appID, appKey}
}

// Name returns the name of the source
func (g *api) Name() string {
	return Name
}

// Define takes a word string and returns a dictionary source.Result
func (g *api) Define(word string) (source.Result, error) {
	// Prepare our URL
	requestURL, err := constructRequestURL(word)

	if nil != err {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)

	if nil != err {
		return nil, err
	}

	httpRequest.Header.Set(httpRequestAcceptHeaderName, jsonMIMEType)
	httpRequest.Header.Set(httpRequestAppIDHeaderName, g.appID)
	httpRequest.Header.Set(httpRequestAppKeyHeaderName, g.appKey)

	httpResponse, err := g.httpClient.Do(httpRequest)

	if nil != err {
		return nil, err
	}

	defer httpResponse.Body.Close()

	if http.StatusNotFound == httpResponse.StatusCode {
		return nil, &source.EmptyResultError{Word: word}
	}

	if http.StatusForbidden == httpResponse.StatusCode {
		return nil, &source.AuthenticationError{}
	}

	if err = source.ValidateHTTPResponse(httpResponse, validMIMETypes, nil); nil != err {
		return nil, err
	}

	body, err := ioutil.ReadAll(httpResponse.Body)

	if nil != err {
		return nil, err
	}

	var result apiResult

	if err = json.Unmarshal(body, &result); nil != err {
		return nil, err
	}

	if len(result.Results) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return source.ValidateAndReturnResult(result.toResult())
}

// toResult converts the proprietary API result to a generic source.Result
func (r apiResult) toResult() source.Result {
	mainResult := r.Results[0]

	entries := make([]interface{}, len(mainResult.LexicalEntries))

	for i, lexicalEntry := range mainResult.LexicalEntries {
		entry := oxfordEntry{}

		for _, pronunciation := range lexicalEntry.Pronunciations {
			if strings.EqualFold(phoneticNotationIPAIdentifier, pronunciation.PhoneticNotation) {
				entry.PronunciationVal = pronunciation.PhoneticSpelling
			}
		}

		entry.WordVal = lexicalEntry.Text
		entry.CategoryVal = lexicalEntry.LexicalCategory.Text

		for _, subEntry := range lexicalEntry.Entries {
			entry.EtymologyVals = append(entry.EtymologyVals, subEntry.Etymologies...)

			for _, sense := range subEntry.Senses {
				senseVal := sense.toSenseValue()

				// Only go one level deep of sub-senses
				for _, subSense := range sense.Subsenses {
					senseVal.SubsenseVals = append(senseVal.SubsenseVals, subSense.toSenseValue())
				}

				entry.SenseVals = append(entry.SenseVals, senseVal)
			}
		}

		entries[i] = entry
	}

	return source.ResultValue{
		Head:      mainResult.Word,
		Lang:      mainResult.Language,
		EntryVals: entries,
	}
}

// toSenseValue converts the proprietary API sense to a source.SenseValue
func (s apiSense) toSenseValue() source.SenseValue {
	examples := make([]string, len(s.Examples))
	notes := make([]string, len(s.Notes))

	for i, example := range s.Examples {
		examples[i] = example.Text
	}

	for i, note := range s.Notes {
		notes[i] = note.Text
	}

	return source.SenseValue{
		DefinitionVals: s.Definitions,
		ExampleVals:    examples,
		NoteVals:       notes,
	}
}

//constructRequestURL constructs the request URL for the Oxford API
func constructRequestURL(word string) (*url.URL, error) {
	return url.Parse(entriesURLString + defaultEntriesSourceLanguage + "/" + word)
}

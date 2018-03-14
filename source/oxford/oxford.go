// Package oxford provides a dictionary source via the Oxford Dictionaries API
//
// Copyright © 2018 Trevor N. Suarez (Rican7)
package oxford

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Rican7/define/source"
)

const (
	// baseURLString is the base URL for all Oxford API interactions
	baseURLString = "https://od-api.oxforddictionaries.com/api/v1/"

	entriesURLString = baseURLString + "entries/"
)

// apiURL is the URL instance used for Oxford API calls
var apiURL *url.URL

// api is a struct containing a configured HTTP client for Oxford API operations
type api struct {
	httpClient *http.Client
	appID      string
	appKey     string
}

// apiResult is a struct that defines the data structure for Oxford API results
type apiResult struct {
	Metadata struct {
		Provider string
	}
	Results []struct {
		ID             string
		Language       string
		LexicalEntries []struct {
			DerivativeOf []struct {
				Domains   []string
				ID        string
				Language  string
				Regions   []string
				Registers []string
				Text      string
			}
			Derivatives []struct {
				Domains   []string
				ID        string
				Language  string
				Regions   []string
				Registers []string
				Text      string
			}
			Entries []struct {
				Etymologies         []string
				GrammaticalFeatures []struct {
					Text string
					Type string
				}
				HomographNumber string
				Notes           []struct {
					ID   string
					Text string
					Type string
				}
				Pronunciations []struct {
					AudioFile        string
					Dialects         []string
					PhoneticNotation string
					PhoneticSpelling string
					Regions          []string
				}
				Senses []struct {
					CrossReferenceMarkers []string
					CrossReferences       []struct {
						ID   string
						Text string
						Type string
					}
					Definitions []string
					Domains     []string
					Examples    []struct {
						Definitions []string
						Domains     []string
						Notes       []struct {
							ID   string
							Text string
							Type string
						}
						Regions      []string
						Registers    []string
						SenseIds     []string
						Text         string
						Translations []struct {
							Domains             []string
							GrammaticalFeatures []struct {
								Text string
								Type string
							}
							Language string
							Notes    []struct {
								ID   string
								Text string
								Type string
							}
							Regions   []string
							Registers []string
							Text      string
						}
					}
					ID    string
					Notes []struct {
						ID   string
						Text string
						Type string
					}
					Pronunciations []struct {
						AudioFile        string
						Dialects         []string
						PhoneticNotation string
						PhoneticSpelling string
						Regions          []string
					}
					Regions   []string
					Registers []string
					Subsenses []struct {
					}
					Translations []struct {
						Domains             []string
						GrammaticalFeatures []struct {
							Text string
							Type string
						}
						Language string
						Notes    []struct {
							ID   string
							Text string
							Type string
						}
						Regions   []string
						Registers []string
						Text      string
					}
					VariantForms []struct {
						Regions []string
						Text    string
					}
				}
				VariantForms []struct {
					Regions []string
					Text    string
				}
			}
			GrammaticalFeatures []struct {
				Text string
				Type string
			}
			Language        string
			LexicalCategory string
			Notes           []struct {
				ID   string
				Text string
				Type string
			}
			Pronunciations []struct {
				AudioFile        string
				Dialects         []string
				PhoneticNotation string
				PhoneticSpelling string
				Regions          []string
			}
			Text         string
			VariantForms []struct {
				Regions []string
				Text    string
			}
		}
		Pronunciations []struct {
			AudioFile        string
			Dialects         []string
			PhoneticNotation string
			PhoneticSpelling string
			Regions          []string
		}
		Type string
		Word string
	}
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

// Define takes a word string and returns a dictionary source.Result
func (g *api) Define(word string) (source.Result, error) {
	// Prepare our URL
	requestURL, err := url.Parse(entriesURLString + "en/" + word)

	if nil != err {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)

	if nil != err {
		return nil, err
	}

	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("app_id", g.appID)
	httpRequest.Header.Set("app_key", g.appKey)

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
	err = json.Unmarshal(body, &result)

	return result.toResult(), err
}

// toResult converts the proprietary API result to a generic source.Result
func (r apiResult) toResult() source.Result {
	mainResult := r.Results[0] // TODO: Handle empty results as an error in the `Define` method

	entries := make([]interface{}, len(mainResult.LexicalEntries))

	for i, lexicalEntry := range mainResult.LexicalEntries {
		entry := oxfordEntry{}

		for _, pronunciation := range lexicalEntry.Pronunciations {
			// TODO: Make a constant
			if strings.EqualFold("IPA", pronunciation.PhoneticNotation) {
				entry.PronunciationVal = pronunciation.PhoneticSpelling
			}
		}

		entry.WordVal = lexicalEntry.Text
		entry.CategoryVal = lexicalEntry.LexicalCategory

		for _, subEntry := range lexicalEntry.Entries {
			entry.EtymologyVals = append(entry.EtymologyVals, subEntry.Etymologies...)

			for _, sense := range subEntry.Senses {
				examples := make([]string, len(sense.Examples))
				notes := make([]string, len(sense.Notes))

				for i, example := range sense.Examples {
					examples[i] = example.Text
				}

				for i, note := range sense.Notes {
					notes[i] = note.Text
				}

				entry.SenseVals = append(
					entry.SenseVals,
					source.SenseValue{
						DefinitionVals: sense.Definitions,
						ExampleVals:    examples,
						NoteVals:       notes,
					},
				)
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

package oxford

// apiResponse defines the data structure for an Oxford API response
type apiResponse struct {
	Metadata struct {
	} `json:"metadata"`
	Results []apiResult `json:"results"`
}

// apiResult defines the data structure for an Oxford API result
type apiResult struct {
	ID             string             `json:"id"`
	Language       string             `json:"language"`
	LexicalEntries []apiLexicalEntry  `json:"lexicalEntries"`
	Pronunciations []apiPronunciation `json:"pronunciations"`
	Type           string             `json:"type"`
	Word           string             `json:"word"`
}

// apiLexicalEntry defines the data structure for an Oxford API lexical entry
type apiLexicalEntry struct {
	Compounds    []apiWordReference `json:"compounds"`
	DerivativeOf []apiWordReference `json:"derivativeOf"`
	Derivatives  []apiWordReference `json:"derivatives"`
	Entries      []struct {
		CrossReferenceMarkers []string           `json:"crossReferenceMarkers"`
		CrossReferences       []apiTypedIDText   `json:"crossReferences"`
		Etymologies           []string           `json:"etymologies"`
		GrammaticalFeatures   []apiTypedIDText   `json:"grammaticalFeatures"`
		HomographNumber       string             `json:"homographNumber"`
		Inflections           []apiInflection    `json:"inflections"`
		Notes                 []apiTypedIDText   `json:"notes"`
		Pronunciations        []apiPronunciation `json:"pronunciations"`
		Senses                []apiSense         `json:"senses"`
		VariantForms          []apiVariantForm   `json:"variantForms"`
	} `json:"entries"`
	GrammaticalFeatures []apiTypedIDText   `json:"grammaticalFeatures"`
	Language            string             `json:"language"`
	LexicalCategory     apiIDText          `json:"lexicalCategory"`
	Notes               []apiTypedIDText   `json:"notes"`
	PhrasalVerbs        []apiWordReference `json:"phrasalVerbs"`
	Phrases             []apiWordReference `json:"phrases"`
	Pronunciations      []apiPronunciation `json:"pronunciations"`
	Root                string             `json:"root"`
	Text                string             `json:"text"`
	VariantForms        []apiVariantForm   `json:"variantForms"`
}

// apiWordReference defines the data structure for an Oxford API word reference
type apiWordReference struct {
	Domains   []apiIDText `json:"domains"`
	ID        string      `json:"id"`
	Language  string      `json:"language"`
	Regions   []apiIDText `json:"regions"`
	Registers []apiIDText `json:"registers"`
	Text      string      `json:"text"`
}

// apiIDText defines the data structure for an Oxford API text with ID
type apiIDText struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// apiTypedIDText defines the data structure for an Oxford API typed, ID'd text
type apiTypedIDText struct {
	apiIDText

	Type string `json:"type"`
}

// apiInflection defines the data structure for an Oxford API inflection
type apiInflection struct {
	Domains             []apiIDText        `json:"domains"`
	GrammaticalFeatures []apiTypedIDText   `json:"grammaticalFeatures"`
	InflectedForm       string             `json:"inflectedForm"`
	LexicalCategory     apiIDText          `json:"lexicalCategory"`
	Pronunciations      []apiPronunciation `json:"pronunciations"`
	Regions             []apiIDText        `json:"regions"`
	Registers           []apiIDText        `json:"registers"`
}

// apiVariantForm defines the data structure for an Oxford API variant form
type apiVariantForm struct {
	Domains        []apiIDText        `json:"domains"`
	Notes          []apiTypedIDText   `json:"notes"`
	Pronunciations []apiPronunciation `json:"pronunciations"`
	Regions        []apiIDText        `json:"regions"`
	Registers      []apiIDText        `json:"registers"`
	Text           string             `json:"text"`
}

// apiSense defines the data structure for an Oxford API "sense"
type apiSense struct {
	Antonyms      []apiWordReference `json:"antonyms"`
	Constructions []struct {
		Domains   []apiIDText      `json:"domains"`
		Examples  [][]string       `json:"examples"`
		Notes     []apiTypedIDText `json:"notes"`
		Regions   []apiIDText      `json:"regions"`
		Registers []apiIDText      `json:"registers"`
		Text      string           `json:"text"`
	} `json:"constructions"`
	CrossReferenceMarkers []string         `json:"crossReferenceMarkers"`
	CrossReferences       []apiTypedIDText `json:"crossReferences"`
	Definitions           []string         `json:"definitions"`
	DomainClasses         []apiIDText      `json:"domainClasses"`
	Domains               []apiIDText      `json:"domains"`
	Etymologies           []string         `json:"etymologies"`
	Examples              []struct {
		Definitions []string         `json:"definitions"`
		Domains     []apiIDText      `json:"domains"`
		Notes       []apiTypedIDText `json:"notes"`
		Regions     []apiIDText      `json:"regions"`
		Registers   []apiIDText      `json:"registers"`
		SenseIds    []string         `json:"senseIds"`
		Text        string           `json:"text"`
	} `json:"examples"`
	ID               string             `json:"id"`
	Inflections      []apiInflection    `json:"inflections"`
	Notes            []apiTypedIDText   `json:"notes"`
	Pronunciations   []apiPronunciation `json:"pronunciations"`
	Regions          []apiIDText        `json:"regions"`
	Registers        []apiIDText        `json:"registers"`
	SemanticClasses  []apiIDText        `json:"semanticClasses"`
	ShortDefinitions []string           `json:"shortDefinitions"`
	Subsenses        []apiSense         `json:"subsenses"`
	Synonyms         []apiWordReference `json:"synonyms"`
	ThesaurusLinks   []struct {
		EntryID string `json:"entry_id"`
		SenseID string `json:"sense_id"`
	} `json:"thesaurusLinks"`
	VariantForms []apiVariantForm `json:"variantForms"`
}

// apiPronunciation defines the data structure for an Oxford API "pronunciation"
type apiPronunciation struct {
	AudioFile        string      `json:"audioFile"`
	Dialects         []string    `json:"dialects"`
	PhoneticNotation string      `json:"phoneticNotation"`
	PhoneticSpelling string      `json:"phoneticSpelling"`
	Regions          []apiIDText `json:"regions"`
	Registers        []apiIDText `json:"registers"`
}

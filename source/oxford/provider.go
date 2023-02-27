// Copyright © 2018 Trevor N. Suarez (Rican7)

package oxford

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	flag "github.com/ogier/pflag"

	"github.com/Rican7/define/registry"
	"github.com/Rican7/define/source"
)

// RequiredConfigError represents an error when a required configuration key is
// missing or invalid.
type RequiredConfigError struct {
	Key string
}

type config struct {
	AppID  string
	AppKey string
}

type provider struct{}

// JSONKey defines the JSON key used for the provider
const JSONKey = "OxfordDictionary"

func init() {
	registry.Register(registry.RegisterFunc(register))
}

func register(flags *flag.FlagSet) (registry.SourceProvider, registry.Configuration) {
	return &provider{}, initConfig(flags)
}

func initConfig(flags *flag.FlagSet) *config {
	conf := &config{}

	// Define our flags
	flags.StringVar(&conf.AppID, "oxford-dictionary-app-id", "", fmt.Sprintf("The app ID for the %s", Name))
	flags.StringVar(&conf.AppKey, "oxford-dictionary-app-key", "", fmt.Sprintf("The app key for the %s", Name))

	return conf
}

func (e *RequiredConfigError) Error() string {
	return fmt.Sprintf("required configuration key %q is missing", e.Key)
}

func (c *config) JSONKey() string {
	return JSONKey
}

// UnmarshalJSON defines how the configuration should be JSON unmarshalled.
func (c *config) UnmarshalJSON(data []byte) error {
	// Alias our type so that we can unmarshal as usual
	type alias config
	copy := &alias{}

	// Unmarshal into our copy
	err := json.Unmarshal(data, copy)
	if err != nil {
		return err
	}

	if c.AppID == "" {
		c.AppID = copy.AppID
	}

	if c.AppKey == "" {
		c.AppKey = copy.AppKey
	}

	return nil
}

func (c *config) Finalize() {
	if c.AppID == "" {
		c.AppID = os.Getenv("OXFORD_DICTIONARY_APP_ID")
	}

	if c.AppKey == "" {
		c.AppKey = os.Getenv("OXFORD_DICTIONARY_APP_KEY")
	}
}

func (p *provider) Name() string {
	return Name
}

func (p *provider) Provide(conf registry.Configuration) (source.Source, error) {
	config := conf.(*config)

	if config.AppID == "" {
		return nil, &RequiredConfigError{Key: "AppID"}
	}

	if config.AppKey == "" {
		return nil, &RequiredConfigError{Key: "AppKey"}
	}

	return New(http.Client{}, config.AppID, config.AppKey), nil
}

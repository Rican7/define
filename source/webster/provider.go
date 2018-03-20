// Copyright © 2018 Trevor N. Suarez (Rican7)

package webster

import (
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
	AppKey string
}

type provider struct{}

// JSONKey defines the JSON key used for the provider
const JSONKey = "MerriamWebsterDictionary"

func init() {
	registry.Register(Name, registry.RegisterFunc(register))
}

func register(flags *flag.FlagSet) (registry.SourceProvider, registry.Configuration) {
	return &provider{}, initConfig(flags)
}

func initConfig(flags *flag.FlagSet) *config {
	conf := &config{}

	// Define our flags
	flags.StringVar(&conf.AppKey, "meriam-webster-dictionary-app-key", "", fmt.Sprintf("The app key for the %s", Name))

	// Attempt to get our values from environment variables
	conf.AppKey = os.Getenv("MERIAM_WEBSTER_DICTIONARY_APP_KEY")

	return conf
}

func (e *RequiredConfigError) Error() string {
	return fmt.Sprintf("Required configuration key %q is missing", e.Key)
}

func (c *config) JSONKey() string {
	return JSONKey
}

func (p *provider) Provide(conf registry.Configuration) (source.Source, error) {
	config := conf.(*config)

	if "" == config.AppKey {
		return nil, &RequiredConfigError{Key: "AppKey"}
	}

	return New(http.Client{}, config.AppKey), nil
}

# CHANGELOG

## vNext TODO

### Features

- Improved internal tooling for more strict linting
- Updated dependencies
- Added support for loading config files in multiple locations, based on the OS/environment
- Added support for the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) for config locations in *nix/Unix-like environments
- Tweaked configuration fallback-order/precedence to use the more common ordering of:
    1. Command line flags
    2. Environment variables
    3. A configuration file
    4. Default values
- Added a `--debug-config` flag to print config file loading status and the paths that are searched for config files

### Bug fixes

- Updated the install instructions in the README
- Made some refactors to satisfy more strict linting
- Fixed the automatic inflection correction fallback searching in the Oxford source to prevent searching for no word (an empty query, which causes an invalid request error)


## 0.3.0

### Features

- Refactored internal data structures
- Added support for showing a source of an example quote
- Added the optional capability for sources to expose a word search
- Added automatic word search result display as a fallback for when define results come back empty
- Added automatic inflection correction fallback to the Oxford source, to align behaviors closer to the other sources (and to improve the user experience)
- Added define result sorting: Now if a direct match is found, it'll be returned first in the list of results
- Added handling of categorical sense data

### Bug fixes

- Fixed the source ambiguity of error messages


## 0.2.0

### Features

- New support for multiple word results for a definition lookup
- Refactored internal data structures
- Improved continuous integration caching
- Updated the Merriam-Webster Dictionary API source to use the new V3 API
- Added a new source: "Free Dictionary API"

### Bug fixes

- Fixed handling of text token cleaning in the Webster source
- Fixed filering of results from the Webster source
- Removed the no-longer functioning "Glosbe API" source

### Security

- Updated dependencies for security issues
    - CVE-2022-41723: https://github.com/golang/go/issues/57855


## 0.1.3

### Features

- Updated the source to use Go 1.20
- Updated all of the go module dependencies
- Internal tooling is now versioned
- Improved internal tooling for more strict linting
- Switched continuous integration testing to use GitHub Actions
- Updated the Oxford Dictionaries API source to use the new V2 API

### Bug fixes

- Updated the install instructions in the README
- Made some refactors to satisfy more strict linting
- Removed use of deprecated functions in favor of their replacements
- Removed now unsupported build targets


## 0.1.2

### Features

- Migrated dependency/package management from `dep` to Go Modules

### Bug fixes

- Improved the cleaning/sanitizing of results from the Glosbe source
- Improved some wording in the README


## 0.1.1

### Features

- A lot of code refactoring, cleanup, and restructuring into more internal packages to prevent shared state and to separate units
- Switched the provider logic so that the "preferred source" configuration allows for fallback, if the "preferred" source can't be provided
- New configuration setting for specifying an exact source that MUST be able to provided or else an error is thrown
- Source listing (via `--list-sources`) now lists both the source names AND their "key" names, so that they can more easily be discovered for configuration
- The Oxford Dictionary API is now the default preferred source, as it provides the highest quality results
- Better handling of errors, especially when reading configurations
- New `--no-config-file` flag to disable the loading of a configuration file

### Bug fixes

- Fixed misspellings of "Merriam"
- Improved the cleaning/sanitizing of results from the Merriam-Webster and Glosbe sources


## 0.1.0

Initial release!

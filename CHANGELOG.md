# CHANGELOG

## vNext TODO

### Features

- New support for multiple word results for a definition lookup
- Refactored internal data structures

### Bug fixes

- 


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

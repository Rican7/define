# CHANGELOG

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

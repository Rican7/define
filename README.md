# define

[![Build Status](https://github.com/Rican7/define/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/Rican7/define/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Rican7/define)](https://goreportcard.com/report/github.com/Rican7/define)
[![Latest Stable Version](https://img.shields.io/github/release/Rican7/define.svg?style=flat)](https://github.com/Rican7/define/releases)

A command-line dictionary (thesaurus) app, with access to multiple sources, written in Go.

<p align="center">
    <img width="822" alt="screen shot 2018-03-21 at 8 51 54 pm" src="https://user-images.githubusercontent.com/742384/37749239-b1b2804e-2d4c-11e8-9e20-f14d1431bbaf.png">
</p>


## Install

Pre-compiled binaries are available on the [releases](https://github.com/Rican7/define/releases) page.

If you have a working Go environment, you can install via `go install`:

```shell
go install github.com/Rican7/define@latest
```


## Configuration

The **define** app allows configuration through multiple means. You can either set configuration via:

- Command line flags (good for one-off use)
- A configuration file (good for your "dotfiles")
- Environment variables (especially useful for API keys)


### Command line flags

The list of command line flags is easily discovered via the `--help` flag. Any passed command line flag will take precedence over any other configuration mechanism.

### Configuration file

A configuration file can be stored at `~/.define.conf.json` and **define** will automatically load the values specified there.

To print the default values of the configuration, simply use the `--print-config` flag. This can also be used to initialize a configuration file, for example:

```shell
define --print-config > ~/.define.conf.json
```

### Environment variables

Some configuration values can also be specified via environment variables. This is especially useful for API keys of different sources.

The following environment variables are read by **define**'s sources:

- `MERRIAM_WEBSTER_DICTIONARY_APP_KEY`
- `OXFORD_DICTIONARY_APP_ID`
- `OXFORD_DICTIONARY_APP_KEY`


## Sources

The **define** app has access to multiple sources, however some of them require user-specific API keys, due to usage limitations.

A preferred source can be specified with the command line flag `--preferred-source="..."` or in a configuration file. For more information, see the section on [Configuration](#configuration).

### Obtaining API keys

The following are links to register for API keys for the different sources:

- [Merriam-Webster's Dictionary API](https://www.dictionaryapi.com/register/index.htm)
- [Oxford Dictionaries API](https://developer.oxforddictionaries.com/?tag=#plans)

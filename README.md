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

1. Command line flags (good for one-off use)
2. Environment variables (good for API keys)
3. A configuration file (good for your "dotfiles")

When multiple means of configuration are used, the values will take precedence in the aforementioned priority.


### Command line flags

The list of command line flags is easily discovered via the `--help` flag. Any passed command line flag will take precedence over any other configuration mechanism.

### Environment variables

Some configuration values can also be specified via environment variables. This is especially useful for API keys of different sources.

The following environment variables are read by **define**'s sources:

- `MERRIAM_WEBSTER_DICTIONARY_APP_KEY`
- `OXFORD_DICTIONARY_APP_ID`
- `OXFORD_DICTIONARY_APP_KEY`

### Configuration file

A configuration file can be stored that **define** will automatically load the values from.

The path of the configuration file to load can be specified via the `--config-file` flag. If no config file path is specified, **define** will search for a config file in your OS's standard config directory paths. While these paths are OS-specific, there are two locations that are searched for that are shared among all platforms:

1. `$XDG_CONFIG_HOME/define/config.json` (This is only searched for when the `$XDG_CONFIG_HOME` env variable is set)
2. `~/.define.conf.json` (Where `~` is equal to your `$HOME` or user directory for your OS)

To see which config file has been loaded, and to check what paths are searched for config files, use the `--debug-config` flag.

To print the default values of the configuration, simply use the `--print-config` flag. This can also be used to initialize a configuration file, for example:

```shell
define --print-config > ~/.define.conf.json
```


## Sources

The **define** app has access to multiple sources, however some of them require user-specific API keys, due to usage limitations.

A preferred source can be specified with the command line flag `--preferred-source="..."` or in a configuration file. For more information, see the section on [Configuration](#configuration).

### Obtaining API keys

The following are links to register for API keys for the different sources:

- [Merriam-Webster's Dictionary API](https://www.dictionaryapi.com/register/index.htm)
- [Oxford Dictionaries API](https://developer.oxforddictionaries.com/?tag=#plans)

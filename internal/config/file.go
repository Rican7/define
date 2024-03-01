package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	xdgBaseName              = "define"
	defaultXDGConfigFileName = "config.json"

	oldDefaultConfigFilePath = "~/.define.conf.json"
)

var (
	userHomeDirPath   string
	userConfigDirPath string
)

// userHomeDir returns the user's home directory, caching the value upon first
// calculation, and without worrying about errors about detection.
func userHomeDir() string {
	if userHomeDirPath == "" {
		// Ignore errors here. We only need the value, even if it's empty.
		userHomeDirPath, _ = os.UserHomeDir()
	}

	return userHomeDirPath
}

// userConfigDir returns the user's config directory, caching the value upon
// first calculation, and without worrying about errors about detection.
func userConfigDir() string {
	if userConfigDirPath == "" {
		// Ignore errors here. We only need the value, even if it's empty.
		userConfigDirPath, _ = os.UserConfigDir()
	}

	return userConfigDirPath
}

// tryExpandUserPath takes a path and expands it if it's user home prefixed (~).
// If the path isn't user home prefixed, then the original path is returned.
func tryExpandUserPath(path string) string {
	if len(path) > 1 && path[0] == '~' && path[1] == filepath.Separator {
		path = filepath.Join(userHomeDir(), path[1:])
	}

	return path
}

// tryExpandConfigPath takes a path and expands it under the user's config home
// path, and the app's config home path within. If those paths aren't known,
// the path is returned empty.
func tryExpandConfigPath(path string) string {
	userConfPath := userConfigDir()

	if userConfPath == "" {
		// If we have no known config path for the user, then we don't have a
		// path we can use.
		return ""
	}

	return filepath.Join(userConfPath, xdgBaseName, path)
}

// findConfigFile attempts to find the current environment user's config file,
// by scanning possible known locations. It returns the path to the config file,
// if any was found.
func findConfigFile() string {
	searchFilePaths := []string{
		tryExpandConfigPath(defaultXDGConfigFileName),
		tryExpandUserPath(oldDefaultConfigFilePath),
	}

	for _, filePath := range searchFilePaths {
		// Check if the file exists
		if _, err := os.Stat(filePath); !errors.Is(err, fs.ErrNotExist) {
			// Return the first file that exists
			// (if there are problems reading the file, we'll handle later)
			return filePath
		}
	}

	return ""
}

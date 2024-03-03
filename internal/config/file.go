package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	xdgBaseName              = "define"
	defaultXDGConfigFileName = "config.json"

	oldDefaultConfigFilePath = "~/.define.conf.json"
)

var (
	userHomeDirPath string
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

// tryExpandUserPath takes a path and expands it if it's user home prefixed (~).
// If the path isn't user home prefixed, then the original path is returned.
func tryExpandUserPath(path string) string {
	if len(path) > 1 && path[0] == '~' && path[1] == filepath.Separator {
		path = filepath.Join(userHomeDir(), path[1:])
	}

	return path
}

// findConfigFile attempts to find the current environment user's config file,
// by scanning possible known locations. It returns the path to the config file,
// if any was found.
func findConfigFile() string {
	defaultXDGConfigRelPath := filepath.Join(xdgBaseName, defaultXDGConfigFileName)

	filePath, err := xdg.SearchConfigFile(defaultXDGConfigRelPath)
	if filePath != "" && err == nil {
		// We found a config! Return it's path.
		return filePath
	}

	oldDefaultConfigFullPath := tryExpandUserPath(oldDefaultConfigFilePath)

	// Check if a file exists at the old default path
	_, err = os.Stat(oldDefaultConfigFullPath)
	if err == nil || errors.Is(err, fs.ErrExist) {
		// Return the file path if it exists
		// (if there are problems reading the file, we'll handle later)
		return oldDefaultConfigFullPath
	}

	return ""
}

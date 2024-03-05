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
	for _, filePath := range FilePaths() {
		// Check if the file exists
		_, err := os.Stat(filePath)
		if err == nil || errors.Is(err, fs.ErrExist) {
			// Return the file path if it exists
			// (if there are problems reading the file, we'll handle later)
			return filePath
		}
	}

	return ""
}

// FilePaths returns the paths of config files that may be searched for in the
// current environment.
//
// This is useful for self-documentation, to provide clarity to users for where
// their config file may be loaded from.
func FilePaths() []string {
	// Length of filePaths is the XDG config home, plus config home, plus the
	// old default file path.
	filePathsLen := len(xdg.ConfigDirs) + 2
	filePaths := make([]string, 0, filePathsLen)

	defaultXDGConfigRelPath := filepath.Join(xdgBaseName, defaultXDGConfigFileName)

	// First we try the user's XDG config home
	filePaths = append(filePaths, filepath.Join(xdg.ConfigHome, defaultXDGConfigRelPath))

	// Then we fall back to the old default path
	filePaths = append(filePaths, tryExpandUserPath(oldDefaultConfigFilePath))

	// Finally, we defer to the XDG config dirs (as those are likely global)
	for _, configDir := range xdg.ConfigDirs {
		filePaths = append(filePaths, filepath.Join(configDir, defaultXDGConfigRelPath))
	}

	return filePaths
}

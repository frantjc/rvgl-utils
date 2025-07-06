package rvglutils

import (
	"os"
	"path/filepath"
)

var (
	SystemPrefPath = func() string {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		return filepath.Join(home, "AppData", "Roaming", "RVGL")
	}()
	DefaultPrefPathList = SystemPrefPath
)

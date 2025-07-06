package rvglutils

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	FlatpakPrefPath = func() string {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		return filepath.Join(home, ".var", "app", "org.rvgl.rvmm/data/rvmm/save")
	}()
	SystemPrefPath = func() string {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		return filepath.Join(home, ".local", "share", "RVGL")
	}()
	DefaultPrefPathList = func() string {
		return strings.Join([]string{SystemPrefPath, FlatpakPrefPath}, string(os.PathListSeparator))
	}()
)

package testdata

import (
	_ "embed"
)

var (
	//go:embed session.csv
	SessionCSV []byte
	//go:embed rvgl.ini
	RVGLINI []byte
	//go:embed profile.ini
	ProfileINI []byte
)

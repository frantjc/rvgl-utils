package testdata

import (
	_ "embed"
)

var (
	//go:embed session.csv
	SessionCSV []byte
)

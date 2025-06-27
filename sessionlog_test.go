package rvglutils_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	rvglutils "github.com/frantjc/rvgl-utils"
	"github.com/frantjc/rvgl-utils/testdata"
)

func TestDecodeSessionCSV(t *testing.T) {
	_, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}
}

func TestResolveSessionCSVByTime(t *testing.T) {
	var (
		tmp      = t.TempDir()
		baseName = "session_1970-01-01_00-00-00.csv"
	)

	if err := os.WriteFile(filepath.Join(tmp, baseName), testdata.SessionCSV, 0644); err != nil {
		t.Fatalf("write %q to %q: %v", baseName, tmp, err)
	}

	name, err := rvglutils.ResolveSessionCSV(&rvglutils.ResolveSessionCSVOpts{
		Path: tmp,
	})
	if err != nil {
		t.Fatalf("resolve %q from %q: %v", baseName, tmp, err)
	}

	file, err := os.Open(name)
	if err != nil {
		t.Fatalf("open %q: %v", name, err)
	}
	t.Cleanup(func() {
		_ = file.Close()
	})

	_, err = rvglutils.DecodeSessionCSV(file)
	if err != nil {
		t.Fatalf("decode %q: %v", filepath.Join(tmp, baseName), err)
	}
}

func TestResolveSessionCSVByName(t *testing.T) {
	var (
		tmp      = t.TempDir()
		baseName = "session_1970-01-01_00-00-00.csv"
	)

	if err := os.WriteFile(filepath.Join(tmp, baseName), testdata.SessionCSV, 0644); err != nil {
		t.Fatalf("write %q to %q: %v", baseName, tmp, err)
	}

	name, err := rvglutils.ResolveSessionCSV(&rvglutils.ResolveSessionCSVOpts{
		Name: baseName,
		Path: tmp,
	})
	if err != nil {
		t.Fatalf("resolve %q from %q: %v", baseName, tmp, err)
	}

	file, err := os.Open(name)
	if err != nil {
		t.Fatalf("open %q: %v", name, err)
	}
	t.Cleanup(func() {
		_ = file.Close()
	})

	_, err = rvglutils.DecodeSessionCSV(file)
	if err != nil {
		t.Fatalf("decode %q: %v", filepath.Join(tmp, baseName), err)
	}
}

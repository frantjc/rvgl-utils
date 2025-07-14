package rvglutils_test

import (
	"os"
	"path/filepath"
	"testing"

	rvglutils "github.com/frantjc/rvgl-utils"
	"github.com/frantjc/rvgl-utils/testdata"
)

func TestResolveAndDecodeSettingsINI(t *testing.T) {
	var (
		tmp      = t.TempDir()
		baseName = "rvgl.ini"
	)

	if err := os.WriteFile(filepath.Join(tmp, baseName), testdata.RVGLINI, 0644); err != nil {
		t.Fatalf("write %q to %q: %v", baseName, tmp, err)
	}

	name, err := rvglutils.ResolveSettingsINI(&rvglutils.ResolveSettingsINIOpts{
		PathList: tmp,
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

	_, err = rvglutils.DecodeSettingsINI(file)
	if err != nil {
		t.Fatalf("decode %q: %v", filepath.Join(tmp, baseName), err)
	}

}

func TestResolveAndDecodeProfileSettingsINI(t *testing.T) {
	var (
		profile  = "frantjc"
		tmp      = t.TempDir()
		tmp2     = filepath.Join(tmp, profile)
		baseName = "profile.ini"
	)

	if err := os.Mkdir(tmp2, 0755); err != nil {
		t.Fatalf("mkdir %q: %v", tmp, err)
	}

	if err := os.WriteFile(filepath.Join(tmp2, baseName), testdata.ProfileINI, 0644); err != nil {
		t.Fatalf("write %q to %q: %v", baseName, tmp2, err)
	}

	name, err := rvglutils.ResolveSettingsINI(&rvglutils.ResolveSettingsINIOpts{
		Profile:  "frantjc",
		PathList: tmp,
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

	_, err = rvglutils.DecodeProfileSettingsINI(file)
	if err != nil {
		t.Fatalf("decode %q: %v", filepath.Join(tmp, baseName), err)
	}
}

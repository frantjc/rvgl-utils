package main

import (
	"context"
	"fmt"

	"github.com/frantjc/rvgl-utils/.dagger/internal/dagger"
)

// +check
func (m *RvglUtilsDev) IsFmted(ctx context.Context) error {
	if empty, err := m.Fmt(ctx).IsEmpty(ctx); err != nil {
		return err
	} else if !empty {
		return fmt.Errorf("source is not formatted (run `dagger call fmt`)")
	}

	return nil
}

// +check
func (m *RvglUtilsDev) TestsPass(
	ctx context.Context,
	// +optional
	githubRepo string,
	// +optional
	githubToken *dagger.Secret,
) error {
	oci := []string{}
	if githubToken != nil && githubRepo != "" {
		oci = append(oci, fmt.Sprintf("ghcr.io/%s/charts/test", githubRepo))
	}

	test, err := m.Test(ctx, oci, githubToken)
	if err != nil {
		return err
	}

	if _, err = test.CombinedOutput(ctx); err != nil {
		return err
	}

	return nil
}

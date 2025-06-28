package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/frantjc/rvgl-utils/command"
	xos "github.com/frantjc/x/os"
)

func main() {
	var (
		cmd       = command.NewRVGLSM()
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	)

	cmd.Version = SemVer()

	err := cmd.ExecuteContext(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			err = nil
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	stop()
	xos.ExitFromError(err)
}

package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	rvglutils "github.com/frantjc/rvgl-utils"
	_ "github.com/frantjc/rvgl-utils/sinks/discord"
	"github.com/frantjc/rvgl-utils/sinks/stdout"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func updateSession(ctx context.Context, sink rvglutils.Sink, sessionCSV string, opts ...rvglutils.UpdateSessionOpt) error {
	file, err := os.Open(sessionCSV)
	if err != nil {
		return fmt.Errorf("open %q: %w", sessionCSV, err)
	}
	defer file.Close() //nolint:errcheck

	session, err := rvglutils.DecodeSessionCSV(file)
	if err != nil {
		return fmt.Errorf("decode %q: %w", sessionCSV, err)
	}

	if len(session.Races) == 0 {
		return nil
	}

	if err = sink.UpdateSession(ctx, session, opts...); err != nil {
		return fmt.Errorf("update session: %w", err)
	}

	return nil
}

// NewRVGLSM returns the command for `rvglsm`.
func NewRVGLSM() *cobra.Command {
	var (
		resolveSessionCSVOpts = &rvglutils.ResolveSessionCSVOpts{}
		scoreSessionOpts      = &rvglutils.ScoreSessionOpts{}
		sinkURL               string
		cmd                   = &cobra.Command{
			Use:           "rvglsm",
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				sessionCSV, err := rvglutils.ResolveSessionCSV(resolveSessionCSVOpts)
				if err != nil {
					return fmt.Errorf("resolve session .csv: %w", err)
				}

				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "resolved session %q\n", sessionCSV)

				var sink rvglutils.Sink = &stdout.Sink{Writer: cmd.OutOrStdout()}

				if sinkURL != "" {
					sink, err = rvglutils.OpenSink(cmd.Context(), sinkURL)
					if err != nil {
						return fmt.Errorf("open sink: %w", err)
					}

					_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "opened sink")
				}

				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					return fmt.Errorf("init file watcher: %w", err)
				}
				defer watcher.Close() //nolint:errcheck

				if err := watcher.Add(filepath.Dir(sessionCSV)); err != nil {
					return fmt.Errorf("watch %q: %w", sessionCSV, err)
				}

				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "began watch on %q\n", sessionCSV)

				go func() {
					for err := range watcher.Errors {
						_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
					}
				}()

				go func() {
					for event := range watcher.Events {
						if event.Name == sessionCSV {
							if err := updateSession(cmd.Context(), sink, sessionCSV, &rvglutils.UpdateSessionOpts{ScoreSessionOpts: scoreSessionOpts}); err != nil {
								watcher.Errors <- err
							}
						}
					}
				}()

				if err := updateSession(cmd.Context(), sink, sessionCSV); err != nil {
					return err
				}
				defer updateSession(context.WithoutCancel(cmd.Context()), sink, sessionCSV, &rvglutils.UpdateSessionOpts{Final: true, ScoreSessionOpts: scoreSessionOpts}) //nolint:errcheck

				<-cmd.Context().Done()
				return cmd.Context().Err()
			},
		}
	)

	cmd.Flags().BoolP("help", "h", false, "Help for "+cmd.Name())
	cmd.Flags().Bool("version", false, "Version for "+cmd.Name())
	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")

	cmd.Flags().StringVarP(&sinkURL, "sink", "s", "", "URL of the sink to send scores to (e.g. discord://{webhook_token}@{webhook_id}[/{message_id}])")
	cmd.Flags().StringVar(&resolveSessionCSVOpts.Name, "session", "", "Name of the session to resolve instead of using the latest one")
	cmd.Flags().BoolVar(&scoreSessionOpts.IncludeAI, "include-ai", false, "Score AI players")
	cmd.Flags().CountVarP(&scoreSessionOpts.ExcludeRaces, "exclude", "x", "Number of races at the beginning of the session to exclude")
	cmd.Flags().StringToIntVarP(&scoreSessionOpts.Handicap, "handicap", "H", nil, "Handicap to apply")

	return cmd
}

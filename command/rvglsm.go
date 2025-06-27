package command

import (
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
				sessionCsv, err := rvglutils.ResolveSessionCSV(resolveSessionCSVOpts)
				if err != nil {
					return fmt.Errorf("resolve session .csv: %w", err)
				}

				fmt.Fprintf(cmd.ErrOrStderr(), "resolved session %q\n", sessionCsv)

				var sink rvglutils.Sink = &stdout.Sink{Writer: cmd.OutOrStdout()}

				if sinkURL != "" {
					sink, err = rvglutils.OpenSink(cmd.Context(), sinkURL)
					if err != nil {
						return fmt.Errorf("open sink: %w", err)
					}

					fmt.Fprintln(cmd.ErrOrStderr(), "opened sink")
				}

				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					return fmt.Errorf("init file watcher: %w", err)
				}
				defer watcher.Close()

				if err := watcher.Add(filepath.Dir(sessionCsv)); err != nil {
					return fmt.Errorf("watch %q: %w", sessionCsv, err)
				}

				fmt.Fprintf(cmd.ErrOrStderr(), "began watch on %q\n", sessionCsv)

				go func() {
					for err := range watcher.Errors {
						fmt.Fprintf(cmd.ErrOrStderr(), "watch: %v\n", err)
					}
				}()

				go func() {
					for event := range watcher.Events {
						if event.Name == sessionCsv {
							file, err := os.Open(sessionCsv)
							if err != nil {
								watcher.Errors <- fmt.Errorf("open %q: %w", sessionCsv, err)
								continue
							}
							defer file.Close()

							session, err := rvglutils.DecodeSessionCSV(file)
							if err != nil {
								watcher.Errors <- fmt.Errorf("decode %q: %w", sessionCsv, err)
								continue
							}

							if err = sink.UpdateScore(cmd.Context(), session, rvglutils.ScoreSession(session, scoreSessionOpts)); err != nil {
								watcher.Errors <- fmt.Errorf("update score: %w", err)
								continue
							}
						}
					}
				}()

				file, err := os.Open(sessionCsv)
				if err != nil {
					watcher.Errors <- fmt.Errorf("open %q: %w", sessionCsv, err)
				}
				defer file.Close()

				session, err := rvglutils.DecodeSessionCSV(file)
				if err != nil {
					watcher.Errors <- fmt.Errorf("decode %q: %w", sessionCsv, err)
				}

				if err := file.Close(); err != nil {
					return err
				}

				if err = sink.UpdateScore(cmd.Context(), session, rvglutils.ScoreSession(session, scoreSessionOpts)); err != nil {
					watcher.Errors <- fmt.Errorf("update score: %w", err)
				}

				<-cmd.Context().Done()
				return cmd.Context().Err()
			},
		}
	)

	cmd.Flags().BoolP("help", "h", false, "Help for "+cmd.Name())
	cmd.Flags().Bool("version", false, "Version for "+cmd.Name())
	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")

	cmd.Flags().StringVarP(&sinkURL, "sink", "s", "", "URL of the sink to send scores to (e.g. discord://{token}@{channel_id})")
	cmd.Flags().BoolVar(&scoreSessionOpts.IncludeAI, "include-ai", false, "Send AI scores to the sink")
	cmd.Flags().StringVar(&resolveSessionCSVOpts.Name, "session", "", "Name of the session to resolve instead of using the latest one")

	return cmd
}

package command

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	rvglutils "github.com/frantjc/rvgl-utils"
	"github.com/frantjc/rvgl-utils/sinks/discord"
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

func init() {
	rvglutils.RegisterSink(&httpSinkOpener{}, "http", "https")
}

type httpSinkOpener struct{}

// Open implements rvglutils.SinkOpener.
func (h *httpSinkOpener) Open(ctx context.Context, u *url.URL) (rvglutils.Sink, error) {
	switch u.Scheme {
	case "http", "https":
	default:
		return nil, fmt.Errorf(`invalid scheme %q, expected "http" or "https"`, u.Scheme)
	}

	switch u.Hostname() {
	case "discordapp.com":
	default:
		return nil, fmt.Errorf(`invalid host %q, expected "discordapp.com"`, u.Hostname())
	}

	var (
		matches   = regexp.MustCompile(`^/api/webhooks/([^/]+)/([^/]+)(?:/messages/([^/]+))?$`).FindStringSubmatch(u.Path)
		messageID string
	)
	switch len(matches) {
	case 4:
		messageID = matches[3]
	case 3:
	default:
		return nil, fmt.Errorf("invalid discord webhook URL path: %q", u.Path)
	}

	return &discord.Sink{
		WebhookID: matches[1],
		Token:     matches[2],
		MessageID: messageID,
	}, nil
}

// NewRVGLSM returns the command for `rvglsm`.
func NewRVGLSM() *cobra.Command {
	var (
		prefPath              string
		resolveSessionCSVOpts = &rvglutils.ResolveSessionCSVOpts{}
		scoreSessionOpts      = &rvglutils.ScoreSessionOpts{}
		laps                  int
		sinkURL               string
		cmd                   = &cobra.Command{
			Use:           "rvglsm",
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				if laps > 0 {
					resolvSettingsINIOpts := &rvglutils.ResolveSettingsINIOpts{PathList: resolveSessionCSVOpts.PathList}

					rvglINI, err := rvglutils.ResolveSettingsINI(resolvSettingsINIOpts)
					if err != nil {
						return err
					}

					settingsFile, err := os.Open(rvglINI)
					if err != nil {
						return err
					}
					defer settingsFile.Close()

					settings, err := rvglutils.DecodeSettingsINI(settingsFile)
					if err != nil {
						return err
					}

					profileINI, err := rvglutils.ResolveSettingsINI(resolvSettingsINIOpts, &rvglutils.ResolveSettingsINIOpts{Profile: settings.Misc.DefaultProfile})
					if err != nil {
						return err
					}

					profileSettingsFile, err := os.Open(profileINI)
					if err != nil {
						return err
					}
					defer profileSettingsFile.Close()

					profileSettings, err := rvglutils.DecodeProfileSettingsINI(profileSettingsFile)
					if err != nil {
						return err
					}

					tmpProfileSettingsFile, err := os.Create(fmt.Sprintf("%s.tmp", profileSettingsFile.Name()))
					if err != nil {
						return err
					}
					defer tmpProfileSettingsFile.Close()

					profileSettings.Game.NLaps = laps

					if err := rvglutils.EncodeSettingsINI(tmpProfileSettingsFile, profileSettings); err != nil {
						return err
					}

					return os.Rename(tmpProfileSettingsFile.Name(), profileSettingsFile.Name())
				}

				sessionCSV, err := rvglutils.ResolveSessionCSV(resolveSessionCSVOpts)
				if err != nil {
					return err
				}

				if prefPath != "" {
					resolveSessionCSVOpts.PathList = filepath.Join(prefPath, "profiles")
				}

				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "resolved session %q\n", sessionCSV)

				var (
					ctx                 = cmd.Context()
					sink rvglutils.Sink = &stdout.Sink{Writer: cmd.OutOrStdout()}
				)

				if sinkURL != "" {
					sink, err = rvglutils.OpenSink(ctx, sinkURL)
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
							if err := updateSession(ctx, sink, sessionCSV, &rvglutils.UpdateSessionOpts{ScoreSessionOpts: scoreSessionOpts}); err != nil {
								watcher.Errors <- err
							}
						}
					}
				}()

				if err := updateSession(ctx, sink, sessionCSV, &rvglutils.UpdateSessionOpts{ScoreSessionOpts: scoreSessionOpts}); err != nil {
					return err
				}
				defer updateSession(context.WithoutCancel(ctx), sink, sessionCSV, &rvglutils.UpdateSessionOpts{Final: true, ScoreSessionOpts: scoreSessionOpts}) //nolint:errcheck

				<-ctx.Done()
				return ctx.Err()
			},
		}
	)

	cmd.Flags().BoolP("help", "h", false, "Help for "+cmd.Name())
	cmd.Flags().Bool("version", false, "Version for "+cmd.Name())
	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")

	cmd.Flags().StringVarP(&sinkURL, "sink", "s", "", "URL of the sink to send updates to (e.g. a Discord webhook URL)")
	cmd.Flags().StringVar(&resolveSessionCSVOpts.Name, "session", "", "Name of the session to resolve instead of using the latest one")
	cmd.Flags().BoolVar(&scoreSessionOpts.IncludeAI, "include-ai", false, "Score AI players")
	cmd.Flags().CountVarP(&scoreSessionOpts.ExcludeRaces, "exclude", "x", "Number of races at the beginning of the session to exclude")
	cmd.Flags().StringToIntVarP(&scoreSessionOpts.Handicap, "handicap", "H", nil, "Handicap to apply")
	cmd.Flags().StringVar(&prefPath, "prefpath", "", "RVGL -prefpath to search for the session in")

	cmd.Flags().IntVar(&laps, "laps", 0, "Set NLaps in default profile.ini and exit")

	return cmd
}

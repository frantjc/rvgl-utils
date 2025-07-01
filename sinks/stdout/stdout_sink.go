package stdout

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	unixtable "github.com/frantjc/go-encoding-unixtable"
	rvglutils "github.com/frantjc/rvgl-utils"
)

func init() {
	rvglutils.RegisterSink(&sinkOpener{}, Scheme)
}

const (
	Scheme = "stdout"
)

type Sink struct {
	io.Writer
}

func (s *Sink) UpdateSession(_ context.Context, session *rvglutils.Session, opts ...rvglutils.UpdateSessionOpt) error {
	o := new(rvglutils.UpdateSessionOpts)

	for _, opt := range opts {
		opt.Apply(o)
	}

	return unixtable.NewEncoder(s.Writer).Encode(rvglutils.ScoreSession(session, o.ScoreSessionOpts))
}

type sinkOpener struct{}

func (o *sinkOpener) Open(ctx context.Context, u *url.URL) (rvglutils.Sink, error) {
	if u.Scheme != Scheme {
		return nil, fmt.Errorf("invalid scheme %s, expected %s", u.Scheme, Scheme)
	}

	s := &Sink{Writer: os.Stdout}

	if u.Path != "" && u.Path != "/" {
		var (
			path = u.Path
			err  error
		)

		switch u.Host {
		case "~":
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}

			path = filepath.Join(home, path)
		default:
			path = filepath.Join(u.Host, path)
		}

		s.Writer, err = os.Open(path)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

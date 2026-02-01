package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	rvglutils "github.com/frantjc/rvgl-utils"
)

func init() {
	rvglutils.RegisterSink(&sinkOpener{}, Scheme)
}

const (
	Scheme = "discord"
)

type Sink struct {
	WebhookID  string
	Token      string
	MessageID  string
	HTTPClient *http.Client
}

// UpdateSession implements rvglutils.Sink.
func (s *Sink) UpdateSession(ctx context.Context, session *rvglutils.Session, opts ...rvglutils.UpdateSessionOpt) error {
	var (
		o       = new(rvglutils.UpdateSessionOpts)
		content = strings.Builder{}
	)

	for _, opt := range opts {
		opt.Apply(o)
	}

	if _, err := content.WriteString(fmt.Sprintf("%s, %s, hosted by %s on %s\n", session.Version, session.Mode, session.Host, session.Date.Format(time.RFC3339))); err != nil {
		return err
	}

	scores := rvglutils.ScoreSession(session, o.ScoreSessionOpts)

	for i, score := range scores {
		format := "%s: %g\n"

		if o.Final && i == 0 {
			format = "**WINNER! %s**: %g\n"
		}

		if _, err := content.WriteString(fmt.Sprintf(format, score.Player, score.Points)); err != nil {
			return err
		}
	}

	u, err := url.Parse(fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", s.WebhookID, s.Token))
	if err != nil {
		return err
	}

	body, err := json.Marshal(map[string]string{"content": content.String()})
	if err != nil {
		return err
	}

	method := http.MethodPost

	if s.MessageID != "" {
		u = u.JoinPath("messages", s.MessageID)
		method = http.MethodPatch
	} else {
		q := u.Query()
		q.Add("wait", "true")
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if s.HTTPClient == nil {
		s.HTTPClient = http.DefaultClient
	}

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close() //nolint:errcheck

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		data := struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{}

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return err
		}

		return fmt.Errorf("%s %s: %d %s", method, u, data.Code, data.Message)
	} else if s.MessageID == "" {
		data := struct {
			ID string `json:"id"`
		}{}

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return err
		} else if data.ID == "" {
			return fmt.Errorf("missing message ID in successful response")
		}

		s.MessageID = data.ID
	}

	return nil
}

type sinkOpener struct{}

// Open implements rvglutils.SinkOpener.
func (o *sinkOpener) Open(ctx context.Context, u *url.URL) (rvglutils.Sink, error) {
	if u.Scheme != Scheme {
		return nil, fmt.Errorf("invalid scheme %q, expected %q", u.Scheme, Scheme)
	}

	var (
		s = &Sink{
			WebhookID: u.Host,
			MessageID: strings.TrimPrefix(u.Path, "/"),
			Token:     u.User.Username(),
		}
	)

	return s, nil
}

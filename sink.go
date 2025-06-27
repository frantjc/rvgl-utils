package rvglutils

import (
	"context"
	"fmt"
	"net/url"
)

type Sink interface {
	UpdateScore(context.Context, *Session, []Score) error
}

type SinkOpener interface {
	Open(context.Context, *url.URL) (Sink, error)
}

var (
	sinkURLMux = map[string]SinkOpener{}
)

func RegisterSink(o SinkOpener, scheme string, schemes ...string) {
	for _, s := range append(schemes, scheme) {
		if _, ok := sinkURLMux[s]; ok {
			panic("attempt to reregister scheme: " + s)
		}

		sinkURLMux[s] = o
	}
}

func OpenSink(ctx context.Context, s string) (Sink, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	if opener, ok := sinkURLMux[u.Scheme]; ok {
		return opener.Open(ctx, u)
	}

	return nil, fmt.Errorf("no sink opener registered for scheme %q", u.Scheme)
}

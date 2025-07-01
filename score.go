package rvglutils

import (
	"sort"
	"strings"
)

type ScoreSessionOpts struct {
	IncludeAI bool
}

func (o *ScoreSessionOpts) Apply(opts *ScoreSessionOpts) {
	if o != nil {
		if opts != nil {
			o.IncludeAI = opts.IncludeAI
		}
	}
}

type ScoreSessionOpt interface {
	Apply(*ScoreSessionOpts)
}

type Score struct {
	Player string
	Points int
}

func newScoreSessionOpts(opts ...ScoreSessionOpt) *ScoreSessionOpts {
	o := &ScoreSessionOpts{}

	for _, opt := range opts {
		opt.Apply(o)
	}

	return o
}

func ScoreSession(session *Session, opts ...ScoreSessionOpt) []Score {
	if session == nil {
		return []Score{}
	}

	var (
		o   = newScoreSessionOpts(opts...)
		tmp = make(map[string]int)
	)

	for _, race := range session.Races {
		players := len(race.Results)

		for _, result := range race.Results {
			if !o.IncludeAI && (result.Car == result.Player || strings.ToUpper(result.Player) != result.Player) {
				continue
			}

			tmp[result.Player] += 1 + players - result.Position
		}
	}

	var (
		score = make([]Score, len(tmp))
		i     = 0
	)
	for player, points := range tmp {
		score[i] = Score{
			Player: player,
			Points: points,
		}
		i++
	}

	sort.Slice(score, func(i, j int) bool {
		return score[i].Points > score[j].Points
	})

	return score
}

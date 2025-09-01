package rvglutils

import (
	"sort"
	"strings"
)

type ScoreSessionOpts struct {
	IncludeAI          bool
	Interval           int
	ExtraPointsPerRace int
	ExcludeRaces       int
	Handicap           map[string]int
	Multipliers        map[string]float64
}

func (o *ScoreSessionOpts) Apply(opts *ScoreSessionOpts) {
	if o != nil {
		if opts != nil {
			opts.IncludeAI = o.IncludeAI
			if o.ExcludeRaces > 0 {
				opts.ExcludeRaces = o.ExcludeRaces
			}
			if o.Interval > 0 {
				opts.Interval = o.Interval
			}
			if o.Handicap != nil {
				opts.Handicap = o.Handicap
			}
			if o.Multipliers != nil {
				opts.Multipliers = o.Multipliers
			}
		}
	}
}

type ScoreSessionOpt interface {
	Apply(*ScoreSessionOpts)
}

type Score struct {
	Player string
	Points float64
}

func newScoreSessionOpts(opts ...ScoreSessionOpt) *ScoreSessionOpts {
	o := &ScoreSessionOpts{}

	for _, opt := range opts {
		opt.Apply(o)
	}

	return o
}

func ScoreSession(session *Session, opts ...ScoreSessionOpt) []Score {
	if session == nil || len(session.Races) == 0 {
		return []Score{}
	}

	var (
		o        = newScoreSessionOpts(opts...)
		tmp      = make(map[string]float64)
		lenRaces = len(session.Races)
	)
	if o.ExcludeRaces > lenRaces {
		o.ExcludeRaces = lenRaces
	} else if o.ExcludeRaces < 0 {
		o.ExcludeRaces = 0
	}

	for k, v := range o.Handicap {
		tmp[k] = float64(v)
	}

	for _, race := range session.Races[o.ExcludeRaces:] {
		players := len(race.Results)

		for _, result := range race.Results {
			if !o.IncludeAI && (result.Car == result.Player || strings.ToUpper(result.Player) != result.Player) {
				continue
			}

			points := float64(1 + o.ExtraPointsPerRace + players - result.Position)
			if points < 0 {
				points = 0
			}

			if o.Multipliers != nil {
				if multiplier, ok := o.Multipliers[result.Car]; ok {
					points *= multiplier
				}
			}

			tmp[result.Player] += points

			if tmp[result.Player] >= float64(o.Interval) && o.Interval > 0 {
				tmp[result.Player] = 0
			}
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

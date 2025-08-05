package rvglutils_test

import (
	"bytes"
	"testing"

	rvglutils "github.com/frantjc/rvgl-utils"
	"github.com/frantjc/rvgl-utils/testdata"
)

func TestScoreSession(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	var (
		scores    = rvglutils.ScoreSession(session)
		lenScores = len(scores)
	)

	if lenScores == 0 {
		t.Fatal("empty score")
	}

	for i, score := range scores {
		if i < lenScores-1 && score.Points < scores[i+1].Points {
			t.Fatal("scores not sorted by points descending")
		}
	}
}

func TestScoreSessionIncludeAI(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	var (
		scores    = rvglutils.ScoreSession(session, &rvglutils.ScoreSessionOpts{IncludeAI: true})
		lenScores = len(scores)
	)

	if lenScores == 0 {
		t.Fatal("empty score")
	}

	if lenScores == 1 {
		t.Fatal("missing AI scores")
	}
}

func TestScoreSessionExclude(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	var (
		scores    = rvglutils.ScoreSession(session, &rvglutils.ScoreSessionOpts{ExcludeRaces: 1})
		lenScores = len(scores)
	)

	if lenScores == 0 {
		t.Fatal("empty score")
	}

	if scores[0].Player != "FRANTJC" {
		t.Fatalf("unexpected player in 1st: %s", scores[0].Player)
	}

	if scores[0].Points != 35 {
		t.Fatal("unexpected 1st place score:", scores[0].Points)
	}
}

func TestScoreSessionExcludeOutOfBounds(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	scores := rvglutils.ScoreSession(session, &rvglutils.ScoreSessionOpts{ExcludeRaces: 9})

	for _, score := range scores {
		if score.Points > 0 {
			t.Fatal("unexpected 1st place score:", scores[0].Points)
		}
	}
}

func TestScoreSessionHandicap(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	var (
		scores    = rvglutils.ScoreSession(session, &rvglutils.ScoreSessionOpts{Handicap: map[string]int{"FRANTJC": 1}})
		lenScores = len(scores)
	)

	if lenScores == 0 {
		t.Fatal("empty score")
	}

	if scores[0].Player != "FRANTJC" {
		t.Fatal("unexpected player in 1st:", scores[0].Player)
	}

	if scores[0].Points != 48 {
		t.Fatal("unexpected 1st place score:", scores[0].Points)
	}
}

func TestScoreSessionIntervalEqualToRacers(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	var (
		scores    = rvglutils.ScoreSession(session, &rvglutils.ScoreSessionOpts{Interval: 12})
		lenScores = len(scores)
	)

	if lenScores == 0 {
		t.Fatal("empty score")
	}

	if scores[0].Player != "FRANTJC" {
		t.Fatal("unexpected player in 1st:", scores[0].Player)
	}

	if scores[0].Points != 11 {
		t.Fatal("unexpected 1st place score:", scores[0].Points)
	}
}

func TestScoreSessionIntervalOffset(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	var (
		scores    = rvglutils.ScoreSession(session, &rvglutils.ScoreSessionOpts{Interval: 24})
		lenScores = len(scores)
	)

	if lenScores == 0 {
		t.Fatal("empty score")
	}

	if scores[0].Player != "FRANTJC" {
		t.Fatal("unexpected player in 1st:", scores[0].Player)
	}

	if scores[0].Points != 23 {
		t.Fatal("unexpected 1st place score:", scores[0].Points)
	}
}

func TestScoreSessionIgnoreAI(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	var (
		scores    = rvglutils.ScoreSession(session, &rvglutils.ScoreSessionOpts{Interval: 24})
		lenScores = len(scores)
	)

	if lenScores == 0 {
		t.Fatal("empty score")
	}

	if scores[0].Player != "FRANTJC" {
		t.Fatal("unexpected player in 1st:", scores[0].Player)
	}

	if scores[0].Points != 23 {
		t.Fatal("unexpected 1st place score:", scores[0].Points)
	}
}

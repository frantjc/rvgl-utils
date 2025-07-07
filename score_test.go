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
			t.Fatal("scores not sorted by points")
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
		t.Fatal("unexpected player in 1st:", scores[0].Player)
	}

	if scores[0].Points != 36 {
		t.Fatal("unexpected 1st place score")
	}
}

func TestScoreSessionExcludeOutOfBounds(t *testing.T) {
	session, err := rvglutils.DecodeSessionCSV(bytes.NewReader(testdata.SessionCSV))
	if err != nil {
		t.Fatalf("decode testdata/session.csv: %v", err)
	}

	rvglutils.ScoreSession(session, &rvglutils.ScoreSessionOpts{ExcludeRaces: 9})
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

	if scores[0].Points != 49 {
		t.Fatal("unexpected 1st place score")
	}
}

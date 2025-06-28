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
		score     = rvglutils.ScoreSession(session)
		lenScores = len(score)
	)

	if lenScores == 0 {
		t.Fatal("empty score")
	}

	for i, s := range score {
		if i < lenScores-1 && s.Points < score[i+1].Points {
			t.Fatal("scores not sorted by points")
		}
	}
}

package rvglutils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ResolveSessionCSVOpts struct {
	Time time.Time
	Name string
	Path string
}

func (o *ResolveSessionCSVOpts) Apply(opts *ResolveSessionCSVOpts) {
	if o != nil {
		if opts != nil {
			if !o.Time.IsZero() {
				opts.Time = o.Time
			}
			if o.Name != "" {
				opts.Name = o.Name
			}
			if o.Path != "" {
				opts.Path = o.Path
			}
		}
	}
}

type ResolveSessionCSVOpt interface {
	Apply(*ResolveSessionCSVOpts)
}

const (
	EnvVarRVGLDataDir = "RVGL_DATA_DIR"
)

func DataDir() (string, error) {
	if dataDir := os.Getenv(EnvVarRVGLDataDir); dataDir != "" {
		return dataDir, nil
	}

	if runtime.GOOS != "linux" {
		return "", fmt.Errorf("unimplemented for %q, please set %q", runtime.GOOS, EnvVarRVGLDataDir)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".var", "app", "org.rvgl.rvmm", "data", "rvmm"), nil
}

func newResolveSessionCSVOpt(opts ...ResolveSessionCSVOpt) *ResolveSessionCSVOpts {
	o := &ResolveSessionCSVOpts{
		Time: time.Now(),
	}

	for _, opt := range opts {
		opt.Apply(o)
	}

	if o.Path == "" {
		if dataDir, err := DataDir(); err == nil {
			o.Path = filepath.Join(dataDir, "save", "profiles")
		}
	}

	return o
}

func ResolveSessionCSV(opts ...ResolveSessionCSVOpt) (string, error) {
	o := newResolveSessionCSVOpt(opts...)

	if o.Name != "" {
		if filepath.IsAbs(o.Name) {
			return o.Name, nil
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dirs := []string{cwd}

	if o.Path != "" {
		parts := strings.Split(o.Path, string(os.PathListSeparator))
		for _, p := range parts {
			if p != "" {
				dirs = append(dirs, p)
			}
		}
	}

	// If Name is set, search for it on Path.
	if o.Name != "" {
		for _, dir := range dirs {
			name := filepath.Join(dir, o.Name)
			if _, err := os.Stat(name); err == nil {
				return name, nil
			}
		}

		return "", fmt.Errorf("%q not found on path %q", o.Name, o.Path)
	}

	var (
		candidateName string
		candidateDiff time.Duration = 1<<63 - 1 // Max duration.
	)
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()

			if !strings.HasPrefix(name, "session_") || !strings.HasSuffix(name, ".csv") {
				continue
			}

			// Extract time.Time from filename like "session_YYYY-MM-DD_HH-MM-SS.csv".
			value := strings.TrimSuffix(strings.TrimPrefix(name, "session_"), ".csv")

			_time, err := time.Parse("2006-01-02_15-04-05", value)
			if err != nil {
				return "", err
			}

			diff := o.Time.Sub(_time)
			if diff < 0 {
				diff = -diff
			}

			if diff < candidateDiff {
				candidateDiff = diff
				candidateName = filepath.Join(dir, name)
			}
		}
	}

	if candidateName != "" {
		return candidateName, nil
	}

	return "", fmt.Errorf("not implemented")
}

type Session struct {
	Version string
	Date    time.Time
	Host    string
	Mode    string
	Laps    int
	AI      bool
	Races   []Race
}

type Race struct {
	Track   string
	Results []Result
}

type Result struct {
	Position int
	Player   string
	Car      string
	Time     time.Duration
	BestLap  time.Duration
	Finished bool
	Cheating bool
}

func DecodeSessionCSV(r io.Reader) (*Session, error) {
	var (
		c          = csv.NewReader(r)
		session    Session
		race       *Race
		lenResults int
		resultIdx  int
	)
	c.FieldsPerRecord = -1

	for {
		record, err := c.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if len(record) == 0 {
			continue
		}

		switch record[0] {
		case "Version":
			if len(record) < 2 {
				return nil, fmt.Errorf("not enough fields in version line")
			}

			session.Version = record[1]
		case "Session":
			if len(record) < 6 {
				return nil, fmt.Errorf("not enough fields in session line")
			}

			session.Date, err = time.Parse(time.ANSIC, record[1])
			if err != nil {
				return nil, err
			}

			session.Host = record[2]
			session.Mode = record[3]
			session.Laps, err = strconv.Atoi(record[4])
			if err != nil {
				return nil, err
			}

			session.AI, err = strconv.ParseBool(record[5])
			if err != nil {
				return nil, err
			}
		case "Results":
			if len(record) < 3 {
				return nil, fmt.Errorf("not enough fields in results line")
			}

			if race != nil {
				session.Races = append(session.Races, *race)
			}

			lenResults, err = strconv.Atoi(record[2])
			if err != nil {
				return nil, fmt.Errorf("parse results count: %v", err)
			}

			resultIdx = 0
			race = &Race{
				Track:   record[1],
				Results: make([]Result, lenResults),
			}
		case "#":
			if race == nil {
				return nil, fmt.Errorf("result header found before race")
			}
		default:
			if len(record) < 7 {
				return nil, fmt.Errorf("not enough fields in result line")
			}

			if race == nil {
				return nil, fmt.Errorf("result found before race")
			}

			position, err := strconv.Atoi(record[0])
			if err != nil {
				return nil, fmt.Errorf("parse position number: %v", err)
			}

			finished, err := strconv.ParseBool(record[5])
			if err != nil {
				return nil, fmt.Errorf("parse finished bool: %v", err)
			}

			cheating, err := strconv.ParseBool(record[6])
			if err != nil {
				return nil, fmt.Errorf("parse cheating bool: %v", err)
			}

			_time, err := parseTime(record[3])
			if err != nil {
				return nil, fmt.Errorf("parse time duration: %v", err)
			}

			bestLap, err := parseTime(record[4])
			if err != nil {
				return nil, fmt.Errorf("parse best lap duration: %v", err)
			}

			race.Results[resultIdx] = Result{
				Position: position,
				Player:   record[1],
				Car:      record[2],
				Time:     _time,
				BestLap:  bestLap,
				Finished: finished,
				Cheating: cheating,
			}
			resultIdx++
		}
	}

	if race != nil {
		session.Races = append(session.Races, *race)
	}

	return &session, nil
}

func parseTime(value string) (time.Duration, error) {
	parts := strings.Split(value, ":")

	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid time format: %s", value)
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("parse minutes: %v", err)
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("parse seconds: %v", err)
	}

	milliseconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("parse milliseconds: %v", err)
	}

	return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(milliseconds)*time.Millisecond, nil
}

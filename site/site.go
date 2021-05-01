package site

import (
	"encoding/csv"
	"errors"
	"os"
	"time"
)

type Metrics struct {
	Address      string
	ResponseTime time.Duration
	Status       int
	Regexp       string
	RegexpStatus string
}

type Entries []Entry

type Entry struct {
	Address string
	Regexp  string
}

func GetEntries(filename string) (Entries, error) {
	var entries Entries
	f, err := os.Open(filename)
	if err != nil {
		return entries, err
	}
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return entries, err
	}
	for _, line := range lines {
		s := Entry{
			Address: line[0],
			Regexp:  line[1],
		}
		// Only add this if we have an address - skip otherwise.
		if s.Address != "" {
			entries = append(entries, s)
		}
	}
	if len(entries) == 0 {
		return entries, errors.New("no entries - cannot proceed")
	}
	return entries, nil
}

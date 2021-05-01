package site

import "time"

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
